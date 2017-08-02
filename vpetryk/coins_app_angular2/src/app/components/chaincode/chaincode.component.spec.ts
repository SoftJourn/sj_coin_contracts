import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ChaincodeComponent } from './chaincode.component';

describe('ChaincodeComponent', () => {
  let component: ChaincodeComponent;
  let fixture: ComponentFixture<ChaincodeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChaincodeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChaincodeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
