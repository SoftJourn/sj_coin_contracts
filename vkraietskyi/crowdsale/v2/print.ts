"use strict";
export class Print {
    public static printBalance(coinContract: any, address: string, message: string): Promise<void> {
        return new Promise<void>(resolve => {
            coinContract.balanceOf(address, function (error, result) {
                if (error) {
                    console.log(error);
                    throw error;
                }
                console.log(message, result.toNumber());
                console.log('------------------------------------------------------------------------------------------------');
                resolve();
            });
        });
    }

    public static printValueOf(contract: any, variable: string, message: string, type: string): Promise<void> {
        return new Promise<void>(resolve => {
            contract[variable](function (error, result) {
                if (error) {
                    console.log(error);
                    throw error;
                }
                if (type == "number") {
                    console.log(message, result.toNumber());
                    console.log('------------------------------------------------------------------------------------------------');
                    resolve();
                } else if (type == "date") {
                    let datetime = new Date(result.toNumber() * 1000);
                    console.log(message, datetime);
                    console.log('------------------------------------------------------------------------------------------------');
                    resolve();
                } else {
                    console.log(message, result);
                    console.log('------------------------------------------------------------------------------------------------');
                    resolve();
                }
            });
        });
    }

    public static printItemsOf(contract: any, getItemMethod: string, lengthMethod: string, message: string): Promise<void> {
        return new Promise<void>(resolve => {
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
                resolve();
            });
        });
    }

    public static print(message: any): Promise<void> {
        return new Promise<void>(resolve => {
            console.log(message);
            resolve();
        });
    }

    public static printVSSeparator(message: any): Promise<void> {
        return new Promise<void>(resolve => {
            console.log(message);
            console.log('------------------------------------------------------------------------------------------------');
            resolve();
        });
    }
}