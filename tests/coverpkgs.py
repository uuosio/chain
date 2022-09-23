import os

coverpkgs = ['..', '../database', '../eosio']

for f in os.listdir('.'):
    if os.path.isdir(f) and f.startswith('test'):
        coverpkgs.append(f'./{f}')

print(','.join(coverpkgs))
