import signal
import threading
import time
from concurrent.futures import ThreadPoolExecutor, as_completed

import boto3
from botocore.exceptions import NoCredentialsError, PartialCredentialsError

# Global flag for graceful shutdown
shutdown_flag = threading.Event()


def signal_handler(signum, frame):
    print("Termination signal received. Initiating graceful shutdown...")
    shutdown_flag.set()


# Register signal handlers
signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)


def get_accounts():
    # Simulating external service call to get accounts and regions
    return [
        {"account": "12345", "region": "ap-south-1"},
        {"account": "67890", "region": "us-west-2"},
        # Add more account-region pairs as needed
    ]


def get_skip_accounts():
    # Simulating external service call to get accounts to be skipped
    return [
        {"account": "12345", "region": "ap-south-1"}
        # Add more account-region pairs to be skipped as needed
    ]


def assume_role(account_id, role_name):
    sts_client = boto3.client("sts")
    response = sts_client.assume_role(
        RoleArn=f"arn:aws:iam::{account_id}:role/{role_name}",
        RoleSessionName="AssumeRoleSession",
    )
    return response["Credentials"]


def requires_upgrade(current_version, latest_version):
    return current_version != latest_version


def get_latest_ami_version(eks_client, cluster_name, nodegroup_name):
    response = eks_client.describe_nodegroup(
        clusterName=cluster_name, nodegroupName=nodegroup_name
    )
    current_version = response["nodegroup"]["releaseVersion"]
    # Assuming the latest version is fetched from an API or predefined value.
    latest_version = (
        "latest-version-from-api-or-config"  # Replace with actual logic
    )
    return current_version, latest_version


def upgrade_nodegroup_ami_version(eks_client, cluster_name, nodegroup_name):
    try:
        response = eks_client.update_nodegroup_version(
            clusterName=cluster_name,
            nodegroupName=nodegroup_name,
            launchTemplate={"version": "$Latest"},
        )
        return response["update"]["id"]
    except Exception as e:
        print(
            f"Error upgrading nodegroup {nodegroup_name} in cluster {cluster_name}: {e}"
        )


def process_region_for_account(account_id, region, role_name):
    if shutdown_flag.is_set():
        print(
            f"Shutdown flag set, skipping processing for account {account_id} "
            f"in region {region}"
        )
        return []

    credentials = assume_role(account_id, role_name)
    eks_client = boto3.client(
        "eks",
        region_name=region,
        aws_access_key_id=credentials["AccessKeyId"],
        aws_secret_access_key=credentials["SecretAccessKey"],
        aws_session_token=credentials["SessionToken"],
    )
    try:
        clusters = eks_client.list_clusters()["clusters"]
        results = []
        for cluster in clusters:
            if shutdown_flag.is_set():
                print(
                    f"Shutdown flag set, stopping processing for cluster "
                    f"{cluster} in account {account_id} region {region}"
                )
                return []

            nodegroups = eks_client.list_nodegroups(clusterName=cluster)[
                "nodegroups"
            ]
            for nodegroup in nodegroups:
                current_version, latest_version = get_latest_ami_version(
                    eks_client, cluster, nodegroup
                )
                if requires_upgrade(current_version, latest_version):
                    update_id = upgrade_nodegroup_ami_version(
                        eks_client, cluster, nodegroup
                    )
                    if update_id:
                        print(
                            f"Successfully initiated upgrade for nodegroup "
                            f"{nodegroup} in cluster {cluster} in region {region} "
                            f"for account {account_id}"
                        )
                        results.append(
                            (eks_client, cluster, nodegroup, update_id)
                        )
        return results
    except Exception as e:
        print(
            f"Error processing region {region} for account {account_id}: {e}"
        )
        return []


def check_upgrade_status(eks_client, cluster_name, nodegroup_name, update_id):
    try:
        response = eks_client.describe_update(
            name=cluster_name, updateId=update_id, nodegroupName=nodegroup_name
        )
        status = response["update"]["status"]
        if status == "InProgress":
            return False, None
        elif status == "Failed":
            return True, (
                f"Upgrade failed for nodegroup {nodegroup_name} in "
                f"cluster {cluster_name}"
            )
        elif status == "Successful":
            return True, (
                f"Successfully upgraded nodegroup {nodegroup_name} in "
                f"cluster {cluster_name}"
            )
    except Exception as e:
        return True, (
            f"Error checking upgrade status for nodegroup {nodegroup_name} in "
            f"cluster {cluster_name}: {e}"
        )

    return False, None


def main(role_name="OrganizationAccountAccessRole"):
    accounts_to_process = mock_get_accounts()
    accounts_to_skip = mock_get_skip_accounts()

    # Report accounts to be skipped
    for account in accounts_to_skip:
        print(
            f"Skipping account {account['account']} in region {account['region']}"
        )

    # Filter out the accounts to be skipped
    filtered_accounts = [
        account_info
        for account_info in accounts_to_process
        if account_info not in accounts_to_skip
    ]

    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = []

        for account_info in filtered_accounts:
            account_id = account_info["account"]
            region = account_info["region"]
            futures.append(
                executor.submit(
                    process_region_for_account, account_id, region, role_name
                )
            )

        # Continuously check the status of nodegroup updates
        while futures:
            # Check for new results from the processing
            for future in as_completed(futures):
                if shutdown_flag.is_set():
                    print("Shutdown flag set, cancelling remaining futures...")
                    for future in futures:
                        future.cancel()
                    futures.clear()
                    break

                try:
                    results = future.result()
                    for eks_client, cluster, nodegroup, update_id in results:
                        futures.append(
                            executor.submit(
                                check_upgrade_status,
                                eks_client,
                                cluster,
                                nodegroup,
                                update_id,
                            )
                        )
                except Exception as e:
                    print(f"Exception occurred: {e}")

                futures.remove(future)

            time.sleep(5)


if __name__ == "__main__":
    try:
        main()
    except (NoCredentialsError, PartialCredentialsError) as e:
        print(f"Credentials error: {e}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")
