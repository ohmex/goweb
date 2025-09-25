package services

import (
	"fmt"
	"goweb/models"
	"goweb/util"

	"gorm.io/gorm"
)

type DomainService struct {
	DB *gorm.DB
}

func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{DB: db}
}

func (service *DomainService) GetDomainByUUID(domain *models.Domain, id string) error {
	return service.DB.
		Where("uuid = ?", id).
		First(domain).Error
}

// CreateDomain creates a new domain and automatically creates a partition for it
func (service *DomainService) CreateDomain(domain *models.Domain) error {
	// Create the domain first
	if err := service.DB.Create(domain).Error; err != nil {
		return err
	}

	// Create partition for the new domain if partitioning is enabled
	if err := service.createPartitionForDomain(domain); err != nil {
		return fmt.Errorf("failed to create partition for domain %s: %w", domain.Name, err)
	}

	return nil
}

// createPartitionForDomain creates a PostgreSQL/YugabyteDB partition for the given domain
func (service *DomainService) createPartitionForDomain(domain *models.Domain) error {
	// Check if partitioning is enabled
	if !util.IsPartitioningEnabled() {
		return nil // Partitioning not enabled, skip
	}

	// Check if we're using PostgreSQL or YugabyteDB (both support partitioning)
	if !util.IsDatabasePartitioningSupported(service.DB) {
		return nil // Only PostgreSQL and YugabyteDB support partitioning
	}

	domainUUID := domain.UUID.String()
	partitionName := util.GeneratePartitionName(domainUUID)

	// Create the partition SQL
	partitionSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s 
		PARTITION OF posts 
		FOR VALUES IN ('%s');
	`, partitionName, domainUUID)

	// Execute the partition creation
	if err := service.DB.Exec(partitionSQL).Error; err != nil {
		return fmt.Errorf("failed to create partition %s: %w", partitionName, err)
	}

	return nil
}

// DeleteDomain deletes a domain and its associated partition
func (service *DomainService) DeleteDomain(domain *models.Domain) error {
	// Delete the partition first
	if err := service.deletePartitionForDomain(domain); err != nil {
		return fmt.Errorf("failed to delete partition for domain %s: %w", domain.Name, err)
	}

	// Delete the domain
	return service.DB.Delete(domain).Error
}

// deletePartitionForDomain deletes the PostgreSQL/YugabyteDB partition for the given domain
func (service *DomainService) deletePartitionForDomain(domain *models.Domain) error {
	// Check if partitioning is enabled
	if !util.IsPartitioningEnabled() {
		return nil // Partitioning not enabled, skip
	}

	// Check if we're using PostgreSQL or YugabyteDB (both support partitioning)
	if !util.IsDatabasePartitioningSupported(service.DB) {
		return nil // Only PostgreSQL and YugabyteDB support partitioning
	}

	domainUUID := domain.UUID.String()
	partitionName := util.GeneratePartitionName(domainUUID)

	// Drop the partition SQL
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s;", partitionName)

	// Execute the partition deletion
	if err := service.DB.Exec(dropSQL).Error; err != nil {
		return fmt.Errorf("failed to drop partition %s: %w", partitionName, err)
	}

	return nil
}

// CreatePartitionsForExistingDomains creates partitions for all existing domains
// This is useful for migration scenarios
func (service *DomainService) CreatePartitionsForExistingDomains() error {
	// Check if partitioning is enabled and supported
	if !util.IsPartitioningEnabled() {
		return nil // Partitioning not enabled, skip
	}

	// Check if we're using PostgreSQL or YugabyteDB (both support partitioning)
	if !util.IsDatabasePartitioningSupported(service.DB) {
		return nil // Only PostgreSQL and YugabyteDB support partitioning
	}

	var domains []models.Domain

	if err := service.DB.Find(&domains).Error; err != nil {
		return fmt.Errorf("failed to fetch existing domains: %w", err)
	}

	for _, domain := range domains {
		if err := service.createPartitionForDomain(&domain); err != nil {
			return fmt.Errorf("failed to create partition for domain %s: %w", domain.Name, err)
		}
	}

	return nil
}
