# How to run Quartermaster on a Kubernetes cluster

<!-- TODO Introduction -->

1. Launch the DGraph database:
    ```
    kubectl apply -k dgraph
    ```

1. Launch Quartermaster:
    ```
    kubectl apply -k qmstr
    ```

1. Wait for the building process to be over:
    ```
    kubectl logs --follow $(kubectl get pods --selector job-name=qmstr -o=name) qmstr-client
    ```

1. Forward two local ports to the following two ports on the DGraph Pod:
    ```
    kubectl port-forward dgraph-0 8000:8000
    ```
    ```
    kubectl port-forward dgraph-0 8080:8080
    ```

1. Open http://localhost:8000/?latest in your browser.

1. Click on "Continue":
    <p align="center">
        <img src="doc/img/dgraph_login.png" alt="DGraph login page" width="75%"/>
    </p>

1. Navigate to the "Console" page
    <!-- TODO -->
    <!-- <p align="center">
        <img src="doc/img/dgraph_console.png" alt="DGraph console page"/>
    </p> -->

1. You should now be able to query the database
    <!-- TODO -->
    <!-- ```graphql
    {
        ...
    }
    ``` -->
