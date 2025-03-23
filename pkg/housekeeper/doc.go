// Package housekeeper provides functionality for managing and maintaining
// database tables with specific conditions and retention policies.
//
// The package defines a Table structure to represent a database table along
// with its associated metadata, such as timestamp and date columns, a condition
// for filtering rows, and a minimum retention period in days.
//
// It also provides a thread-safe mechanism to register and manage tables
// using a global map protected by a mutex.
package housekeeper
