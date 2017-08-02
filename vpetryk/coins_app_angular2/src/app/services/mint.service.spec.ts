import { TestBed, inject } from '@angular/core/testing';

import { MintService } from './mint.service';

describe('MintService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [MintService]
    });
  });

  it('should be created', inject([MintService], (service: MintService) => {
    expect(service).toBeTruthy();
  }));
});
