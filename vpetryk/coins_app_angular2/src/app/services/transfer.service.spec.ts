import { TestBed, inject } from '@angular/core/testing';

import { TransferService } from './transfer.service';

describe('TransferService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [TransferService]
    });
  });

  it('should be created', inject([TransferService], (service: TransferService) => {
    expect(service).toBeTruthy();
  }));
});
