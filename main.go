package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// function for getting data from api.binance.com and presenting it in console in table form
// assetsName is a name of assets for which the orders need to be retrieved (e.g. BTCUSDT or BNBUSDT, see https://binance-docs.github.io/apidocs/spot/en/#public-api-definitions for info)
func getOrders(assetsName string) {
	//get data from api.binance.com
	resp, err := http.Get("https://api.binance.com/api/v3/depth?symbol=" + assetsName)
	if err != nil {
		fmt.Println("Error in http.get:", err)
	}

	//present data as []byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in reading body data:", err)
	}
	defer resp.Body.Close()

	//parsing data in body (JSON) to var orders
	var orders = map[string][][]string{}
	err1 := json.Unmarshal(body, &orders)
	if err1 != nil {
		// fmt.Println("error in json.Unmarshal:", err1) //removed error printing, error in converting int of orders ID (which useless in the program) to string, TODO appropriate interface for retrieving JSON data
	}

	//converting string slices of BIDS and ASKS to float slices to make calculations of total sums later
	var bids [][]float64 //slice for saving bids
	var asks [][]float64 //slice for saving asks
	var price float64    //temp var for converting price of an order from string type to float
	var amount float64   //temp var for converting ammount of an order from string type to float
	//saving only 15 last bids and ask (because they are the newest ones), so that the last data in string slice (var orders) becomes the first data in float slices (var bids and asks)
	for i := len(orders["bids"]) - 1; i > len(orders["bids"])-16; i-- {
		if price, err = strconv.ParseFloat(orders["bids"][i][0], 64); err != nil {
			fmt.Println("Error in parsing:", err)
		}
		if amount, err = strconv.ParseFloat(orders["bids"][i][1], 64); err != nil {
			fmt.Println("Error in parsing:", err)
		}
		bids = append(bids, []float64{price, amount})
	}
	for i := len(orders["asks"]) - 1; i > len(orders["asks"])-16; i-- {
		if price, err = strconv.ParseFloat(orders["asks"][i][0], 64); err != nil {
			fmt.Println("Error in parsing:", err)
		}
		if amount, err = strconv.ParseFloat(orders["asks"][i][1], 64); err != nil {
			fmt.Println("Error in parsing:", err)
		}
		asks = append(asks, []float64{price, amount})
	}

	//clear console (!!! work only for Windows!!!), TODO clear console for other platforms
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout //comment this 3 lines if you want to save all printed data in console (for debugging)
	cmd.Run()

	//printing data in console in table form
	//printing headers of table
	fmt.Println("\t\t\t\t\t\t", assetsName)
	fmt.Println("\t\t      BIDS", "\t\t\t\t\t         ASKS")
	fmt.Println("           price     amount       total", "\t\t       price     amount       total")

	//vars for orders sums
	var bidsPriceSum, bidsAmountSum, bidsTotalSum, asksPriceSum, asksAmountSum, asksTotalSum float64

	//printing data line by line and calculate sums
	for i := 0; i < len(bids) || i < len(asks); i++ {
		fmt.Printf("     %12.5f %9.5f %12.5f \t\t %12.5f %9.5f %12.5f \n", bids[i][0], bids[i][1], bids[i][0]*bids[i][1], asks[i][0], asks[i][1], asks[i][0]*asks[i][1])
		bidsPriceSum += bids[i][0]
		bidsAmountSum += bids[i][1]
		bidsTotalSum += bids[i][0] * bids[i][1]
		asksPriceSum += asks[i][0]
		asksAmountSum += asks[i][1]
		asksTotalSum += asks[i][0] * asks[i][1]
	}

	//printing sums of orders
	fmt.Println("")
	fmt.Printf("Sums %12.5f %9.5f %12.5f \t\t %12.5f %9.5f %12.5f \n", bidsPriceSum, bidsAmountSum, bidsTotalSum, asksPriceSum, asksAmountSum, asksTotalSum)
	//printing instruction on how to end the program
	fmt.Println("")
	fmt.Println("Press Ctrl+c to finish execution of the program")
}

func main() {
	//inputting assets name for retrieving orders data
	var assetsName string
	fmt.Println("For which assets do you want to get orders? (e.g. BNBBTC, BTCUSDT...)")
	fmt.Scanf("%s\n", &assetsName)

	//call the function every second in infinite loop
	for {
		getOrders(assetsName)
		time.Sleep(time.Second)
	}
}
