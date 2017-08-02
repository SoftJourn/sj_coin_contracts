import { TestBed, inject } from '@angular/core/testing';

import { ChaincodeService } from './chaincode.service';

describe('ChaincodeService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ChaincodeService]
    });
  });

  it('should be created', inject([ChaincodeService], (service: ChaincodeService) => {
    expect(service).toBeTruthy();
  }));
});
