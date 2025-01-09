// linkedlist/linkedlist_test.go
package stake

import (
	//"main/types"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrentInsertAndRead(t *testing.T) {
	list := NewDoublyLinkedList(10)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Printf("this is %d thread", i)
			list.Insert(id, rand.Intn(100))
		}(i)
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			list.Getlinklist(5)
		}()
	}
	wg.Wait()
}

func TestDoublyLinkedList(t *testing.T) {
	t.Run("Test Insert and Getlinklist", func(t *testing.T) {
		list := NewDoublyLinkedList(5)

		// Test case 1: insert a few elements
		list.Insert(1, 10)
		list.Insert(2, 5)
		list.Insert(3, 15)
		list.Insert(4, 2)
		list.Insert(5, 8)

		expected := []string{"3=15", "1=10", "5=8", "2=5", "4=2"}
		actual := list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 1 Failed. Expected: %v, Got: %v", expected, actual)
		}

		// Test case 2: insert a new node with a larger value than existing nodes
		list.Insert(6, 20)
		expected = []string{"6=20", "3=15", "1=10", "5=8", "2=5"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 2 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 3: insert a new node with a  value  same  existing nodes
		list.Insert(7, 1)

		expected = []string{"6=20", "3=15", "1=10", "5=8", "2=5"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 3 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 4: insert a new node with a smaller value than existing nodes
		list.Insert(7, 12)
		expected = []string{"6=20", "3=15", "7=12", "1=10", "5=8"}

		actual = list.Getlinklist(5)

		if !equal(actual, expected) {
			t.Errorf("Test Case 4 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 5: insert node more than maxSize
		list.Insert(8, 30)
		expected = []string{"8=30", "6=20", "3=15", "7=12", "1=10"}
		actual = list.Getlinklist(5)

		if !equal(actual, expected) {
			t.Errorf("Test Case 5 Failed. Expected: %v, Got: %v", expected, actual)
		}
		list.Insert(9, 18)
		expected = []string{"8=30", "6=20", "9=18", "3=15", "7=12"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 6 Failed. Expected: %v, Got: %v", expected, actual)
		}
		list.Insert(9, 35)
		expected = []string{"9=35", "8=30", "6=20", "3=15", "7=12"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 7 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 8 : insert same value and id , not update
		list.Insert(9, 35)
		expected = []string{"9=35", "8=30", "6=20", "3=15", "7=12"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 8 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 9 : insert node to tail
		list.Insert(10, 1)
		expected = []string{"9=35", "8=30", "6=20", "3=15", "7=12"}
		actual = list.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 9 Failed. Expected: %v, Got: %v", expected, actual)
		}

	})

	t.Run("Test Getlinklist with Limit", func(t *testing.T) {
		list := NewDoublyLinkedList(5)
		list.Insert(1, 10)
		list.Insert(2, 5)
		list.Insert(3, 15)
		list.Insert(4, 2)
		list.Insert(5, 8)
		// Test case 1: get with limit
		expected := []string{"3=15", "1=10", "5=8"}
		actual := list.Getlinklist(3)
		if !equal(actual, expected) {
			t.Errorf("Test Case 1 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 2: get all when limit is more than size
		expected = []string{"3=15", "1=10", "5=8", "2=5", "4=2"}
		actual = list.Getlinklist(6)
		if !equal(actual, expected) {
			t.Errorf("Test Case 2 Failed. Expected: %v, Got: %v", expected, actual)
		}
		// Test case 3: get 0  is nil
		expected = []string{}
		actual = list.Getlinklist(0)
		if !equal(actual, expected) {
			t.Errorf("Test Case 3 Failed. Expected: %v, Got: %v", expected, actual)
		}
		list1 := NewDoublyLinkedList(5)
		// Test case 4: get all nil link
		expected = []string{}
		actual = list1.Getlinklist(5)
		if !equal(actual, expected) {
			t.Errorf("Test Case 4 Failed. Expected: %v, Got: %v", expected, actual)
		}
	})
	t.Run("Test Concurrent Operations", func(t *testing.T) {
		list := NewDoublyLinkedList(10)
		var wg sync.WaitGroup
		numGoroutines := 10
		iterations := 100
		wg.Add(numGoroutines)
		// 模拟并发写入
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					id := rand.Intn(100) // 使用随机ID
					value := rand.Intn(1000)
					list.Insert(id, value)
					// 短暂睡眠，模拟不同goroutine之间的竞争
					time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
				}
			}(i)
		}
		wg.Wait()
		// 模拟并发读取，验证数据是否一致
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				_ = list.Getlinklist(10)
			}()
		}
		wg.Wait()
		// 添加检查点，检查链表长度是否符合预期
		expectedSize := len(list.nodeMap) // 获取map的长度
		actualSize := list.Size
		if actualSize != expectedSize {
			t.Errorf("Concurrent insert  failed. Expected size: %d, Got size: %d", expectedSize, actualSize)

		}
	})
	t.Run("Test Concurrent Read", func(t *testing.T) {
		list := NewDoublyLinkedList(10)
		list.Insert(1, 10)
		list.Insert(2, 20)
		list.Insert(3, 30)
		var wg sync.WaitGroup
		numGoroutines := 10

		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				_ = list.Getlinklist(10)

			}()
		}
		wg.Wait()

	})
	t.Run("Test Concurrent Insert and Update", func(t *testing.T) {
		list := NewDoublyLinkedList(10)
		var wg sync.WaitGroup
		numGoroutines := 10
		iterations := 100
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					id := rand.Intn(10) + 1 // 使用随机ID
					value := rand.Intn(1000)
					list.Insert(id, value)
				}
			}(i)
		}
		wg.Wait()
		actualSize := list.Size
		if actualSize == 0 {
			t.Errorf("TestConcurrentInsertAndUpdate failed, actualSize is 0")
		}
		expectedSize := len(list.nodeMap) // 获取map的长度
		if actualSize != expectedSize {
			t.Errorf("TestConcurrentInsertAndUpdate  failed. Expected size: %d, Got size: %d", expectedSize, actualSize)

		}

	})
	t.Run("Test Concurrent Insert and remove", func(t *testing.T) {
		list := NewDoublyLinkedList(10)
		var wg sync.WaitGroup
		numGoroutines := 10
		iterations := 100
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(goroutineID int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					id := rand.Intn(10) + 1 // 使用随机ID
					value := rand.Intn(1000)
					list.Insert(id, value)
					node := list.findNodeById(id)
					if node != nil {
						list.removeNode(node)
					}
				}
			}(i)
		}
		wg.Wait()
		actualSize := list.Size
		if actualSize != 0 {
			t.Errorf("Test Concurrent Insert and remove  failed, actualSize:%d", actualSize)
		}
	})

}

// Helper function to compare string slices
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func main() {
	list := NewDoublyLinkedList(5)
	list.Insert(1, 10)
	list.Insert(2, 5)
	list.Insert(3, 15)
	list.Insert(4, 2)
	list.Insert(5, 8)
	list.Insert(6, 20)
	list.Insert(7, 12)
	list.Insert(8, 30)
	list.Insert(9, 18)
	list.Insert(9, 35)

	fmt.Println(list.Getlinklist(5)) // Output: [9=35 8=30 6=20 3=15 7=12]
}
