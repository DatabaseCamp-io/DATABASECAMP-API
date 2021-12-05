package services

// interface.go
/**
 * 	This file used to be a interface of services
 */

/**
 * 	 Interface to show function in aws service that others can use
 */
type IAwsService interface {

	/**
	 * Get File link from amazon s3
	 *
	 * @param 	imagekey  object key for getting file link
	 *
	 * @return file link
	 * @return the error of getting file link
	 */
	GetFileLink(imagekey string) (string, error)
}
