package oss

import (
	"gitee.com/wuntsong-auth/backend/src/config"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	errors "github.com/wuntsong-org/wterrors"
)

var IdentityOSS *oss.Client
var IdentityBucket *oss.Bucket

var IdentitySignOSS *oss.Client
var IdentitySignBucket *oss.Bucket

var HeaderOSS *oss.Client
var HeaderBucket *oss.Bucket

var HeaderSignOSS *oss.Client
var HeaderSignBucket *oss.Bucket

var FileOSS *oss.Client
var FileBucket *oss.Bucket

var FileSignOSS *oss.Client
var FileSignBucket *oss.Bucket

var WorkOrderFileOSS *oss.Client
var WorkOrderFileBucket *oss.Bucket

var WorkOrderFileSignOSS *oss.Client
var WorkOrderFileSignBucket *oss.Bucket

var InvoiceOSS *oss.Client
var InvoiceBucket *oss.Bucket

var InvoiceSignOSS *oss.Client
var InvoiceSignBucket *oss.Bucket

func InitOss() errors.WTError {
	if len(config.BackendConfig.Aliyun.AccessKeyId) == 0 {
		return errors.Errorf("aliyun AccessKeyId must be given")
	}

	if len(config.BackendConfig.Aliyun.AccessKeySecret) == 0 {
		return errors.Errorf("aliyun AccessKeySecret must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.Endpoint) == 0 {
		return errors.Errorf("identity oss endpoint must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.BucketName) == 0 {
		return errors.Errorf("identity oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.Sign.Endpoint) == 0 {
		return errors.Errorf("identity sign oss endpoint must be given")
	}

	if len(config.BackendConfig.Aliyun.Identity.Sign.BucketName) == 0 {
		return errors.Errorf("identity sign oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.Endpoint) == 0 {
		return errors.Errorf("header oss endpoint must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.BucketName) == 0 {
		return errors.Errorf("header oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.Header.Sign.Endpoint) == 0 {
		config.BackendConfig.Aliyun.Header.Sign.Endpoint = config.BackendConfig.Aliyun.Header.Endpoint
	}

	if len(config.BackendConfig.Aliyun.Header.Sign.BucketName) == 0 {
		config.BackendConfig.Aliyun.Header.Sign.BucketName = config.BackendConfig.Aliyun.Header.BucketName
	}

	if len(config.BackendConfig.Aliyun.File.Endpoint) == 0 {
		return errors.Errorf("file oss endpoint must be given")
	}

	if len(config.BackendConfig.Aliyun.File.BucketName) == 0 {
		return errors.Errorf("file oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.File.Sign.Endpoint) == 0 {
		config.BackendConfig.Aliyun.File.Sign.Endpoint = config.BackendConfig.Aliyun.File.Endpoint
	}

	if len(config.BackendConfig.Aliyun.File.Sign.BucketName) == 0 {
		config.BackendConfig.Aliyun.File.Sign.BucketName = config.BackendConfig.Aliyun.File.BucketName
	}

	if len(config.BackendConfig.Aliyun.WorkOrder.Endpoint) == 0 {
		return errors.Errorf("work order file oss endpoint must be given")
	}

	if len(config.BackendConfig.Aliyun.WorkOrder.BucketName) == 0 {
		return errors.Errorf("work order file oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.WorkOrder.Sign.Endpoint) == 0 {
		config.BackendConfig.Aliyun.WorkOrder.Sign.Endpoint = config.BackendConfig.Aliyun.WorkOrder.Endpoint
	}

	if len(config.BackendConfig.Aliyun.WorkOrder.Sign.BucketName) == 0 {
		config.BackendConfig.Aliyun.WorkOrder.Sign.BucketName = config.BackendConfig.Aliyun.WorkOrder.BucketName
	}

	if len(config.BackendConfig.Aliyun.Invoice.Endpoint) == 0 {
		return errors.Errorf("invoice file oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.Invoice.BucketName) == 0 {
		return errors.Errorf("invoice file oss bucket name must be given")
	}

	if len(config.BackendConfig.Aliyun.Invoice.Sign.Endpoint) == 0 {
		config.BackendConfig.Aliyun.Invoice.Sign.Endpoint = config.BackendConfig.Aliyun.Invoice.Endpoint
	}

	if len(config.BackendConfig.Aliyun.Invoice.Sign.BucketName) == 0 {
		config.BackendConfig.Aliyun.Invoice.Sign.BucketName = config.BackendConfig.Aliyun.Invoice.BucketName
	}

	var err error
	IdentityOSS, err = oss.New(config.BackendConfig.Aliyun.Identity.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	IdentityBucket, err = IdentityOSS.Bucket(config.BackendConfig.Aliyun.Identity.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	IdentitySignOSS, err = oss.New(config.BackendConfig.Aliyun.Identity.Sign.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	IdentitySignBucket, err = IdentitySignOSS.Bucket(config.BackendConfig.Aliyun.Identity.Sign.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	HeaderOSS, err = oss.New(config.BackendConfig.Aliyun.Header.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	HeaderBucket, err = HeaderOSS.Bucket(config.BackendConfig.Aliyun.Header.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	HeaderSignOSS, err = oss.New(config.BackendConfig.Aliyun.Header.Sign.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	HeaderSignBucket, err = HeaderSignOSS.Bucket(config.BackendConfig.Aliyun.Header.Sign.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	FileOSS, err = oss.New(config.BackendConfig.Aliyun.File.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	FileBucket, err = FileOSS.Bucket(config.BackendConfig.Aliyun.File.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	FileSignOSS, err = oss.New(config.BackendConfig.Aliyun.File.Sign.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	FileSignBucket, err = FileSignOSS.Bucket(config.BackendConfig.Aliyun.File.Sign.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	WorkOrderFileOSS, err = oss.New(config.BackendConfig.Aliyun.WorkOrder.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	WorkOrderFileBucket, err = WorkOrderFileOSS.Bucket(config.BackendConfig.Aliyun.WorkOrder.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	WorkOrderFileSignOSS, err = oss.New(config.BackendConfig.Aliyun.WorkOrder.Sign.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	WorkOrderFileSignBucket, err = WorkOrderFileSignOSS.Bucket(config.BackendConfig.Aliyun.WorkOrder.Sign.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	InvoiceOSS, err = oss.New(config.BackendConfig.Aliyun.Invoice.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	InvoiceBucket, err = InvoiceOSS.Bucket(config.BackendConfig.Aliyun.Invoice.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	InvoiceSignOSS, err = oss.New(config.BackendConfig.Aliyun.Invoice.Sign.Endpoint,
		config.BackendConfig.Aliyun.AccessKeyId,
		config.BackendConfig.Aliyun.AccessKeySecret)
	if err != nil {
		return errors.WarpQuick(err)
	}

	InvoiceSignBucket, err = InvoiceSignOSS.Bucket(config.BackendConfig.Aliyun.Invoice.Sign.BucketName)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
