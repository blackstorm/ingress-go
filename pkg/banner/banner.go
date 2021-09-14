package banner

import "fmt"

const banner = `
.__              ____                                                   ____            
|__|   ____     / ___\  _______    ____     ______   ______            / ___\    ____   
|  |  /    \   / /_/  > \_  __ \ _/ __ \   /  ___/  /  ___/   ______  / /_/  >  /  _ \  
|  | |   |  \  \___  /   |  | \/ \  ___/   \___ \   \___ \   /_____/  \___  /  (  <_> ) 
|__| |___|  / /_____/    |__|     \___  > /____  > /____  >          /_____/    \____/  
          \/                          \/       \/       \/                              
`

func Print(isHiddenBanner *bool) {
	if isHiddenBanner == nil || !*isHiddenBanner {
		fmt.Println(banner)
	}
}
