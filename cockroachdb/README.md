# CockroachDB

## Project Setup

1. Download and extract `cockroach` executable into the `/temp` folder of each node in the cluster.

1. Set the Go variables in the `~/.bashrc` or `~/.bash_profile` as follows:

   ```bash
   export GO_HOME=/temp/cs5424-team-m/go
   export GOPATH=/home/stuproj/cs4224m/go
   export CDB_PATH=/temp/cs5424-team-m/cockroach
   export CDB_SERVER=/temp/cs5424-team-m/cdb-server
   export PATH=$GO_HOME/bin:$CDB_PATH:$CDB_SERVER:$PATH
   ```

## Project Initialization

1. Install Go(go1.13.4)

1. Download the cockroach code into the directory relative to the path mentioned below:

   ```bash
   $ $GOPATH/src/github/com/mychewcents/tpcc-benchmarks
   ```

1. `cd` into the `cockroachdb` directory

   ```bash
   $ cd cockroachdb
   ```

1. Run `bash scripts/build_exec.sh`

   - This script will create all the executable required to perform all the experiments.

1. Update the `configs/prod/setup.json` file to provide the IP addresses of the nodes and their related port numbers. A sample coonfiguration file looks like below:

   ```json
   {
      "data_files_url": "http://www.comp.nus.edu.sg/~cs4224/project-files.zip",
      "working_dir": "/temp/cs5424-team-m",
      "nodes": [
          {
              "id": 1,
              "host": "192.168.48.179",
              "port": 27000,
              "http_addr": "0.0.0.0:40000",
              "username": "root",
              "database": "defaultdb"
          },
          ...
   }
   ```

1. Run the following command to download the dataset required for the experiments.

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 download-dataset
   ```

1. Run the following commands on each node to initialize the cockroach db node files and place the project's data files in the `extern` directory of the node's files.

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 setup-dir
   ```

   NOTE:

   - `node` is the NODE ID of the current node in the cluster and should be between 1 and 5

## Running Experiments

Before every experiment, we just need to run the above `./setupCmd` again to clean the directories. It should be run on all the individual nodes.

1. To start the cockroach db nodes, use the following command:

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 start
   ```

1. To initialize the cockroach cluster, the following command should be run only once for the cluster on any node. Just need to make sure that the `node` value is corresponding to the Node's ID itself.

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 init
   ```

1. To load the data with the initial dataset, one should run the command:

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 load
   ```

   The command will execute the SQL scripts present in the `scripts/raw` directory. As the `load` command only loads the `CSV` files for the project's raw data files, it uses `SQL` queries to create the data required for our experiments to run.

   One can export this data via the command:

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 export
   ```

   NOTE: The export for the non-partitioned tables are on the basis of their names. However, for the partitioned tables, the exports follow the format of `<warehouse id>_<district id>.csv`.

   Once exported, for any following initializations, one can use the `load-csv` command as follows:

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 load-csv
   ```

1. To run the experiments, we need to run the following command on each of the node. We need to make sure that the experiment number provided below is between 5 and 8 as mentioned in the Project's specification document.

   ```bash
   $ ./setupCmd -env=prod -config=configs/prod/setup.json -node=1 -exp=6 run-exp
   ```

   Example: In the above case, the command will launch 4 instances of the client for the particular `node` running on the localhost follwoing the `mod` rule mentioned in the project's specification document.

## Logging

As logs are essential to track progress, one can check the `logs` folder for the related logs for the `setupCmd`, `setupCmd`, and the transactions and SQL files run by them. The naming convention is as follows:

- `setup_*` contains the logs of the `setupCmd` run at the start
- `exp_%d_client_$d_*` contains the logs for the particular experiment and the client's number for that experiment.

NOTE: The last value is the Unix Timestamp to allow for unique logs files everytime a command is executed.

You can run the following command to view the live log updates:

```bash
$ tail -f logs/<log file name>`
```

## Performance Metrics Tracking

One can monitor for the client metrics to trace client progress.

Each client that completes all transactions successfully stores its metrics (the metrics needed for `clients.csv`) for that experiment as a CSV file inside folder inside `results/metrics`. Each client creates a CSV file with format `<experiment number>_<clientNo>.csv`.

You can run `ls results/metrics` to check which all clients have completed their executions.

## Database State

One can get the database state by running the following command:

```bash
$ ./dbstateCmd -env=prod -config=configs/prod/setup.json -node=1
```

This command will fetch the database state information needed and stores it in a CSV file under the folder `results/dbstate`. There will be a file for each experiment with the format `<experiment>.csv`.
