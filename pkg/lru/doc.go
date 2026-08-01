// Package lru provides a generic, thread-safe Least-Recently-Used cache.
//
// The cache evicts the least recently used entry once it reaches its capacity.
// Recency is tracked with a doubly-linked list from the list package: the head
// holds the most-recently-used entry and the tail holds the least-recently-used
// entry. A map provides O(1) lookup from key to its list node and value.
package lru
