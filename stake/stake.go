package stake

import (
	"log"
	"sync"
)

type StakeMap struct {
	StakeMap sync.Map // betOfferId -> stakeMap

}

func NewstakeMap() *StakeMap {
	return &StakeMap{
		StakeMap: sync.Map{},
	}
}
func (sm *StakeMap) Insert(custmerID int, betOfferID int, value int, maxHighStakes int) {

	oldlist, ok := sm.StakeMap.Load(betOfferID) // 使用Load获取，如果不存在则创建
	log.Printf("in post")
	if !ok {
		list := NewDoublyLinkedList(maxHighStakes)
		list.Insert(custmerID, value)
		log.Printf("add in linklist%d ", custmerID)
		sm.StakeMap.Store(betOfferID, list) // 使用 Store 添加数据
		// betMap = list
	} else {
		olist := oldlist.(*DoublyLinkedList) // 断言类型
		olist.Insert(custmerID, value)
		sm.StakeMap.Store(betOfferID, olist) // 使用 Store 添加数据

	}

}
func (sm *StakeMap) GetTop(betOfferID int, maxHighStakes int) ([]string, bool) {

	// stakeMapValue2, ok := sm.StakeMap.Load(9)
	// if !ok {
	// 	log.Printf(" failure")
	// }
	log.Printf(" betid is %d", betOfferID)

	stakeMapValue, ok := sm.StakeMap.Load(betOfferID)
	if !ok {
		return make([]string, 0), false

	}
	log.Printf(" gettop")

	stakeMap := stakeMapValue.(*DoublyLinkedList)
	topStakes := stakeMap.Getlinklist(maxHighStakes)
	return topStakes, true
}
