# Navy Federal Spending Tracker

## Why?

I try to track my expenses each month but sometime it just gets hard to keep up with them.

I created this script to make it easier to get an idea about my expenses each month.

## Usage

The first step is to download the CSV files for the account that you would like to actually get a summary for. Place the CSVs in a separate directory on their own. I suggest naming the CSV in accordance with the account that the file is for. This program for both the transaction data for credit cards and checking accounts. 


```bash
$ git clone https://github.com/Soups71/SpendingTracker.git

$ cd SpendingTracker
$ go build .
$ ./SpendingTracker FinData 4 2024
```

There can be as many CSV files in the folder provided as you would like.


