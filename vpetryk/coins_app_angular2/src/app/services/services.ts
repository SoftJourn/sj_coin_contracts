import {UserService} from './user.service';
import {DataService} from "./data.service";
import {ChannelService} from "./channel.service";
import {ChaincodeService} from "./chaincode.service";
import {MintService} from "./mint.service";
import {TransferService} from "./transfer.service";
import {BalanceService} from "./balance.service";

export let servicesInjectables: Array<any> = [
  ChaincodeService,
  UserService,
  ChannelService,
  DataService,
  MintService,
  TransferService,
  BalanceService
];
