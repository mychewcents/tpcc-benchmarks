import pandas as pd
import numpy as np
import glob

path = 'results/metrics' # use your path
all_files = glob.glob(path + "/*.csv")

li = []

for filename in all_files:
    df = pd.read_csv(filename, index_col=None, header=None)
    li.append(df)

clients = pd.concat(li, axis=0, ignore_index=True)
clients = clients.sort_values(by=[0, 1])
clients.to_csv('results/clients.csv', header=False, index=False)

throughput = clients[[0,4]].groupby(0).agg([np.min, np.mean, np.max])
throughput = throughput.round(2)
throughput.to_csv('results/throughput.csv', header=False)

path = 'results/dbstate' # use your path
all_files = glob.glob(path + "/*.csv")

li = []

for filename in all_files:
    df = pd.read_csv(filename, index_col=None, header=None)
    li.append(df)

db_state = pd.concat(li, axis=0, ignore_index=True)
db_state.to_csv('results/db-state.csv', header=False, index=False)