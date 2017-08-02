import { BrowserModule } from '@angular/platform-browser';
import {NgModule} from '@angular/core';

import { AppComponent } from './components/app/app.component';
import {servicesInjectables} from './services/services';
import { HttpModule } from '@angular/http';
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import { ChannelComponent } from './components/channel/channel.component';
import { ChaincodeComponent } from './components/chaincode/chaincode.component';
import { MintComponent } from './components/mint/mint.component';
import { TransferComponent } from './components/transfer/transfer.component';

@NgModule({
  declarations: [
    AppComponent,
    ChannelComponent,
    ChaincodeComponent,
    MintComponent,
    TransferComponent
  ],
  imports: [
    BrowserModule,
    HttpModule,
    FormsModule,
    ReactiveFormsModule
  ],
  providers: [
    servicesInjectables
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
