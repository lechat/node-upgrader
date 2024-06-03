package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"go.uber.org/zap"
)

type Collector struct {
	logger    *zap.Logger
	config    *Config
	configDir string
}

func NewCollector(config *Config, configDir string, logger *zap.Logger) *Collector {
	return &Collector{
		logger:    logger,
		configDir: configDir,
		config:    config,
	}
}

func (c *Collector) Start() {
	c.logger.Info("Starting collector",
		zap.String("config_dir", c.configDir),
	)

	accounts, err := c.readAccounts()
	if err != nil {
		c.logger.Error("Error reading accounts", zap.Error(err))
		return
	}
	c.logger.Info("Finished reading accounts")

	c.processAccounts(accounts)
}

func (c *Collector) readAccounts() ([]Account, error) {
	accPath := filepath.Join(c.configDir, "accounts.json")
	file, err := os.Open(accPath)
	if err != nil {
		c.logger.Error("Error opening %s file: %s", zap.String("accPath", accPath), zap.Error(err))
		return nil, err
	}
	defer file.Close()

	// Parse the JSON content
	var accounts []Account
	if err := json.NewDecoder(file).Decode(&accounts); err != nil {
		c.logger.Error("Error decoding %s file: %s", zap.String("accPath", accPath), zap.Error(err))
		return nil, err
	}
	return accounts, nil
}

// processAccounts processes accounts in parallel with a limit of maxConcurrency
func (c *Collector) processAccounts(accounts []Account) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, c.config.AccountBatchSize)

	for _, acc := range accounts {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(acc Account) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if err := c.checkEKSClusters(acc); err != nil {
				c.logger.Error("Error checking EKS clusters", zap.String("account_id", acc.Account), zap.Error(err))
			}
		}(acc)
	}

	wg.Wait()
	close(semaphore)
}

// checkEKSClusters checks if the account has EKS clusters and lists node groups
func (c *Collector) checkEKSClusters(account Account) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(account.Region),
	})
	if err != nil {
		c.logger.Error("Error creating AWS session", zap.Error(err))
		return err
	}

	eksSvc := eks.New(sess)

	// List EKS clusters
	clusterResp, err := eksSvc.ListClusters(nil)
	if err != nil {
		c.logger.Error("Error listing EKS clusters", zap.Error(err))
		return err
	}

	for _, clusterName := range clusterResp.Clusters {
		c.logger.Info("Checking EKS cluster", zap.String("cluster_name", *clusterName))
	}

	return nil
}
