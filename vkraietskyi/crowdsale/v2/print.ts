"use strict";
export class Print {
    public static printBalance(coinContract: any, address: string, message: string): void {
        coinContract.balanceOf(address, function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            console.log(message, result.toNumber());
            console.log('------------------------------------------------------------------------------------------------');
        });
    }

    public static printValueOf(contract: any, variable: string, message: string, type: string): void {
        contract[variable](function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            if (type == "number") {
                console.log(message, result.toNumber());
                console.log('------------------------------------------------------------------------------------------------');
            } else if (type == "date") {
                let datetime = new Date(result.toNumber() * 1000);
                console.log(message, datetime);
                console.log('------------------------------------------------------------------------------------------------');
            } else {
                console.log(message, result);
                console.log('------------------------------------------------------------------------------------------------');
            }
        });
    }

    public static printItemsOf(contract: any, getItemMethod: string, lengthMethod: string, message: string): void {
        contract[lengthMethod](function (error, result) {
            if (error) {
                console.log(error);
                throw error;
            }
            Print.printVSSeparator(message);
            for (let i = 0; i < result.toNumber(); i++) {
                contract[getItemMethod](i, function (error, result) {
                    if (error) {
                        console.log(error);
                        throw error;
                    }
                    Print.print(result.toString());
                });
            }
        });
    }

    public static print(message: any): void {
        console.log(message);
    }

    public static printVSSeparator(message: any): void {
        console.log(message);
        console.log('------------------------------------------------------------------------------------------------');
    }
}