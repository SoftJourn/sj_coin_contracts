import { Component, OnInit } from '@angular/core';
import {DataService} from "../../services/data.service";
import {ChannelService} from "../../services/channel.service";

@Component({
  selector: 'app-channel',
  templateUrl: './channel.component.html',
  styleUrls: ['./channel.component.css']
})
export class ChannelComponent implements OnInit {

  constructor(private data:DataService, private channelService:ChannelService) { }

  public isEnabled:boolean = false;

  ngOnInit() {

  }

  getChannels():void {
    if (this.data.user.token) {
      this.isEnabled = true;
      this.channelService.getChannels();
    }
    else {
      this.isEnabled = false;
      console.log("Incorrect user token")
    }
  }

  createChannel():void {
    if (this.data.channel.currentChannel) {
      console.log("Channel already exists")
   }
   else {
      this.channelService.createChannel();
    }
  }

  connectToChannel() {
    if (this.data.channel.isConnected) {
      console.log("already connected")
    }
    else {
      this.channelService.joinChannel();
    }
  }

}
