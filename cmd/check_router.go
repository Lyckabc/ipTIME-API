//go:build ignore

package main

import (
	"fmt"

	"github.com/Lyckabc/ipTIME-API/cmd/routers"
)

func main() {
	router := RouterInfo()
	client := CreateClient()
	routers.Login(client, router)

	fmt.Println("=== 포트 포워딩 목록 ===")
	list := routers.GetPortForwardList(client, router)
	if len(list) == 0 {
		fmt.Println("  (없음)")
	}
	for _, pf := range list {
		fmt.Printf("  %-15s %s  %d→%d  active=%v\n", pf.Name, pf.Target, pf.Src.Start, pf.Dst.Start, pf.Active)
	}

	fmt.Println("=== WOL 목록 ===")
	wols := routers.GetWOLList(client, router)
	if len(wols) == 0 {
		fmt.Println("  (없음)")
	}
	for _, w := range wols {
		fmt.Printf("  %-15s %s\n", w.Name, w.MAC)
	}
}
