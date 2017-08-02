import { Injectable } from '@angular/core';
import { UserModel } from '../models/userModel';
import {ChannelModel} from "../models/channelModel";
import {ChaincodeModel} from "../models/chaincodeModel";

@Injectable()
export class DataService {

  public user: UserModel;
  public channel: ChannelModel;
  public chaincode: ChaincodeModel;

  constructor() {
    this.user = new UserModel();
    this.channel = new ChannelModel();
    this.chaincode = new ChaincodeModel();
  }
}
