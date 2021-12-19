package services

// interface.go
/**
 * 	This file used to be a interface of services
 */

/**
 * 	 Interface to show function in stroage service that others can use
 */
type IStroageService interface {

	/**
	 * Get signed File link from stroage service
	 *
	 * @param 	imagekey  object key for getting file link
	 *
	 * @return file link
	 * @return the error of getting file link
	 */
	GetFileLink(objectName string) (string, error)
}
