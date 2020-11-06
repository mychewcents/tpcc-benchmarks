# Cassandra

## 1. Project Setup

1. Download and extract cassandra into the `/temp` folder of each node in the cluster.
2. Pick a node as a seed node.
3. Update the `conf/cassandra.yaml` file on each node as follows:
```bash
cluster_name: 'ProjectCluster'
num_tokens: 256
seed_provider:
  - class_name: org.apache.cassandra.locator.SimpleSeedProvider
    parameters:
         - seeds:  "<ip_of_seed_node>"
listen_address: if seed node then <ip_of_seed_node> or else blank
rpc_address: 0.0.0.0
endpoint_snitch: GossipingPropertyFileSnitch
read_request_timeout_in_ms: 100000
range_request_timeout_in_ms: 100000
write_request_timeout_in_ms: 100000
counter_write_request_timeout_in_ms: 100000
cas_contention_timeout_in_ms: 100000
truncate_request_timeout_in_ms: 120000
request_timeout_in_ms: 100000
slow_query_log_timeout_in_ms: 50000
```
4. Update the `conf/cassandra-rackdc.properties` file on each node as follows:
```bash
dc=DC1
rack=RAC1
```
5. Delete the `conf/cassandra-topology.properties` file on each node.
6. First start Cassandra on seed node by executing `bin/cassandra`. Then start the other nodes one by one using the same command. 
You can monitor the status of the cluster by using the command `bin/nodetool status`.



## 2. Project Initialization

1. Install Python(3.7.3) and Go(go1.13.4)
2. Download the cassandra code into the home directory and `cd` into the directory.
3. Run `pip install -r requirements.txt`
4. Run `bash scripts/make-initial-data.sh`. 
    - This command will download the data and transaction files for the project. The files will be downloaded to directory `assests/data`.
    - It will then create the processed data which corresponds to our new schema design from the data provided. This is done using a `python` script which can be found in `scripts/python/make-initial-data.py`.  We use the `pandas` library to process the CSV files. This processed data is stored under the directory  `assets/data/processed`.
    - We create this processed data only once.
5. Run `bash scripts/initialize-cassandra.sh`
    - This script will initialize the Cassandra schema and loads initial data.
    - It will initialize the schema using a `cql` script found under `scripts/cql/schema-initialization.cql`.
    - It will then populate the schema with the initial data using a `cql` script which copies the data from the processed CSV files into their respective tables using Cassandra's `COPY` command. This scriot can be found under `scripts/cql/data-initialization.cql`.



## 3. Running Experiments

Need to follow the below steps before running each experiment.

1. `cd` into Cassandra's home directory in each node.
2. Stop Cassandra on each node by running the command `bin/nodetool stopdaemon` from inside of Cassandra's home directory.
3. Delete the folders `data` and `log` inside of Cassandra's home directory  in each node.
   1. Run `rm -rf data`
   2. Run `rm -rf logs`
4. First start Cassandra on seed node by executing `bin/cassandra`. Then start the other nodes one by one using the same command. 
   You can monitor the status of the cluster by using the command `bin/nodetool status`.
5. Now `cd` into the project source code directory.
6. Run `bash scripts/initialize-cassandra.sh` to initialize the data.
7. The Cassandra session configuration such as the IP Address of the nodes and Consistency for each experiment is provided using a config file. You can find the config files under the `configs/prod` folder.
8. To run an experiment you can run the script `bash scripts/run-experiment.sh prod <experimentNo> <hostNo>`.
   - Here `experimentNo` represents the experiment you want to run. Based on this value the script will fetch the configuration from the `configs` folder and runs the required no. of clients in parallel for that experiment on the node.
   - You can assign each node an identifier starting from 1. The `hostNo` represents the node identifier on which the script is being run.
   - You need to run this script on each node with its respective `hostNo`.
9. You can check the log files to trace the clients progress.
   - Each client logs Starting, Stopping and Error log into a log file in directory `log`. The log file will be of format `logs_exp_<experimentNo>_client_<clientNo>`. 
   - You can run `tail -f log/logs_exp_<experimentNo>*` to see the content of all log files for a particular experiment.
   - Each client prints the message `Stopping Client` when it completes its execution successfully. If any error occurs the process will stop immediately and the error is logged in the log file.
10. You can monitor for the client metrics to trace client progress.
    - Each client that completes all transactions successfully stores its metrics (the metrics needed for `clients.csv`) for that experiment  as a CSV file inside folder inside `results/metrics`. Each client creates a CSV file with format `experiment_<experimentNo>_client_<clientNo>.csv`.
    - You can run `ls results/metrics` to check which all clients have completed their executions.
11. After all clients have completed their execution; run `go cmd/dbstate/cassandra-state.go`. This command will fetch the database state information needed and stores it in a CSV file under the folder `results/dbstate`. There will be a file for each experiment with the format `experiment_<experimentNo>.csv`.
12. Run `bash scripts/performance.sh`. This calls a `python` script `scripts/python/performance.py` which uses the `pandas` library to consolidate the clients data into a single file. 
    - It will create a single CSV file `results/clients.csv` from all files under `results/metrics`.
    - It will create a single CSV file `results/throughput.csv` by aggregating the values from all files under `results/metrics`.
    - It will create a single CSV file `results/db-state` from all files under `results/dbstate`.
    - After all experiments are completed the files `results/clients.csv`, `results/throughput.csv` and `results/db-state` will have the data for all the experiments.







