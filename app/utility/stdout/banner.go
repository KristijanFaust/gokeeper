package stdout

import "fmt"

const banner = `
   ____       _  __                         
  / ___| ___ | |/ /___  ___ _ __   ___ _ __ 
 | |  _ / _ \| ' // _ \/ _ \ '_ \ / _ \ '__|
 | |_| | (_) | . \  __/  __/ |_) |  __/ |   
  \____|\___/|_|\_\___|\___| .__/ \___|_|   
                           |_|              
`

func PrintApplicationBanner() {
	fmt.Println(banner)
}
