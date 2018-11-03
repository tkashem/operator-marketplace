package datastore

const RedHatOperatorManifest = `
name: Marketplace Operators
publisher: Marketplace

data:
  customResourceDefinitions: |-
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: etcdbackups.etcd.database.coreos.com
      spec:
        group: etcd.database.coreos.com
        version: v1beta2
        scope: Namespaced
        names:
          kind: EtcdBackup
          listKind: EtcdBackupList
          plural: etcdbackups
          singular: etcdbackup
      
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: etcdclusters.etcd.database.coreos.com
      spec:
        group: etcd.database.coreos.com
        version: v1beta2
        scope: Namespaced
        names:
          plural: etcdclusters
          singular: etcdcluster
          kind: EtcdCluster
          listKind: EtcdClusterList
          shortNames:
            - etcdclus
            - etcd
      
    - apiVersion: apiextensions.k8s.io/v1beta1
      kind: CustomResourceDefinition
      metadata:
        name: etcdrestores.etcd.database.coreos.com
      spec:
        group: etcd.database.coreos.com
        version: v1beta2
        scope: Namespaced
        names:
          kind: EtcdRestore
          listKind: EtcdRestoreList
          plural: etcdrestores
          singular: etcdrestore

  clusterServiceVersions: |-
    - #! validate-crd: ./deploy/chart/templates/03-clusterserviceversion.crd.yaml
      #! parse-kind: ClusterServiceVersion
      apiVersion: operators.coreos.com/v1alpha1
      kind: ClusterServiceVersion
      metadata:
        name: etcdoperator.v0.6.1
        namespace: placeholder
        annotations:		
          tectonic-visibility: ocs	
          packageId: akashem/etcd
          description: Red Hat etcd
          createdAt: 2018/01/01
          certifiedLevel: good
          repository: https://etcd.redhat.com
          containerImage: quay.io/redhat/etcd
          support: Red Hat
          healthIndex: A          	
      spec:
        displayName: etcd
        description: |
          etcd is a distributed key value store that provides a reliable way to store data across a cluster of machines. It’s open-source and available on GitHub. etcd gracefully handles leader elections during network partitions and will tolerate machine failure, including the leader. Your applications can read and write data into etcd.
          A simple use-case is to store database connection details or feature flags within etcd as key value pairs. These values can be watched, allowing your app to reconfigure itself when they change. Advanced uses take advantage of the consistency guarantees to implement database leader elections or do distributed locking across a cluster of workers.
      
          _The etcd Open Cloud Service is Public Alpha. The goal before Beta is to fully implement backup features._
      
          ### Reading and writing to etcd
      
          Communicate with etcd though its command line utility etcdctl or with the API using the automatically generated Kubernetes Service.
      
          [Read the complete guide to using the etcd Open Cloud Service](https://coreos.com/tectonic/docs/latest/alm/etcd-ocs.html)
      
          ### Supported Features
          **High availability**
          Multiple instances of etcd are networked together and secured. Individual failures or networking issues are transparently handled to keep your cluster up and running.
          **Automated updates**
          Rolling out a new etcd version works like all Kubernetes rolling updates. Simply declare the desired version, and the etcd service starts a safe rolling update to the new version automatically.
          **Backups included**
          Coming soon, the ability to schedule backups to happen on or off cluster.
        keywords: ['etcd', 'key value', 'database', 'coreos', 'open source']
        version: 0.6.1
        maturity: alpha
        maintainers:
        - name: CoreOS, Inc
          email: support@coreos.com
      
        provider:
          name: CoreOS, Inc
        labels:
          alm-status-descriptors: etcdoperator.v0.6.1
          alm-owner-etcd: etcdoperator
          operated-by: etcdoperator
        selector:
          matchLabels:
            alm-owner-etcd: etcdoperator
            operated-by: etcdoperator
        links:
        - name: Blog
          url: https://coreos.com/etcd
        - name: Documentation
          url: https://coreos.com/operators/etcd/docs/latest/
        - name: etcd Operator Source Code
          url: https://github.com/coreos/etcd-operator
      
        icon:
        - base64data: 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
          mediatype: image/png
        install:
          strategy: deployment
          spec:
            permissions:
            - serviceAccountName: etcd-operator
              rules:
              - apiGroups:
                - etcd.database.coreos.com
                resources:
                - etcdclusters
                verbs:
                - "*"
              - apiGroups:
                - storage.k8s.io
                resources:
                - storageclasses
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - pods
                - services
                - endpoints
                - persistentvolumeclaims
                - events
                verbs:
                - "*"
              - apiGroups:
                - apps
                resources:
                - deployments
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - secrets
                verbs:
                - get
            deployments:
            - name: etcd-operator
              spec:
                replicas: 1
                selector:
                  matchLabels:
                    name: etcd-operator-alm-owned
                template:
                  metadata:
                    name: etcd-operator-alm-owned
                    labels:
                      name: etcd-operator-alm-owned
                  spec:
                    serviceAccountName: etcd-operator
                    containers:
                    - name: etcd-operator
                      command:
                      - etcd-operator
                      - --create-crd=false
                      image: quay.io/coreos/etcd-operator@sha256:bd944a211eaf8f31da5e6d69e8541e7cada8f16a9f7a5a570b22478997819943
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
        customresourcedefinitions:
          owned:
          - name: etcdclusters.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdCluster
            displayName: etcd Cluster
            description: Represents a cluster of etcd nodes.
            resources:
              - kind: Service
                version: v1
              - kind: Pod
                version: v1
            specDescriptors:
              - description: The desired number of member Pods for the etcd cluster.
                displayName: Size
                path: size
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podCount'
            statusDescriptors:
              - description: The status of each of the member Pods for the etcd cluster.
                displayName: Member Status
                path: members
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podStatuses'
              - description: The service at which the running etcd cluster can be accessed.
                displayName: Service
                path: service
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes:Service'
              - description: The current size of the etcd cluster.
                displayName: Cluster Size
                path: size
              - description: The current version of the etcd cluster.
                displayName: Current Version
                path: currentVersion
              - description: 'The target version of the etcd cluster, after upgrading.'
                displayName: Target Version
                path: targetVersion
              - description: The current status of the etcd cluster.
                displayName: Status
                path: phase
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase'
              - description: Explanation for the current status of the cluster.
                displayName: Status Details
                path: reason
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
      
    - #! validate-crd: ./deploy/chart/templates/03-clusterserviceversion.crd.yaml
      #! parse-kind: ClusterServiceVersion
      apiVersion: operators.coreos.com/v1alpha1
      kind: ClusterServiceVersion
      metadata:
        name: etcdoperator.v0.9.0
        namespace: placeholder
        annotations:
          tectonic-visibility: ocs
          alm-examples: '[{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdCluster","metadata":{"name":"example","namespace":"default"},"spec":{"size":3,"version":"3.2.13"}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdRestore","metadata":{"name":"example-etcd-cluster"},"spec":{"etcdCluster":{"name":"example-etcd-cluster"},"backupStorageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdBackup","metadata":{"name":"example-etcd-cluster-backup"},"spec":{"etcdEndpoints":["<etcd-cluster-endpoints>"],"storageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}}]'
          packageId: akashem/etcd
          description: Red Hat etcd
          createdAt: 2018/01/01
          certifiedLevel: good
          repository: https://etcd.redhat.com
          containerImage: quay.io/redhat/etcd
          support: Red Hat
          healthIndex: A           
      spec:
        displayName: etcd
        description: |
          etcd is a distributed key value store that provides a reliable way to store data across a cluster of machines. It’s open-source and available on GitHub. etcd gracefully handles leader elections during network partitions and will tolerate machine failure, including the leader. Your applications can read and write data into etcd.
          A simple use-case is to store database connection details or feature flags within etcd as key value pairs. These values can be watched, allowing your app to reconfigure itself when they change. Advanced uses take advantage of the consistency guarantees to implement database leader elections or do distributed locking across a cluster of workers.
      
          _The etcd Open Cloud Service is Public Alpha. The goal before Beta is to fully implement backup features._
      
          ### Reading and writing to etcd
      
          Communicate with etcd though its command line utility etcdctl or with the API using the automatically generated Kubernetes Service.
      
          [Read the complete guide to using the etcd Open Cloud Service](https://coreos.com/tectonic/docs/latest/alm/etcd-ocs.html)
      
          ### Supported Features
      
      
          **High availability**
      
      
          Multiple instances of etcd are networked together and secured. Individual failures or networking issues are transparently handled to keep your cluster up and running.
      
      
          **Automated updates**
      
      
          Rolling out a new etcd version works like all Kubernetes rolling updates. Simply declare the desired version, and the etcd service starts a safe rolling update to the new version automatically.
      
      
          **Backups included**
      
      
          Coming soon, the ability to schedule backups to happen on or off cluster.
        keywords: ['etcd', 'key value', 'database', 'coreos', 'open source']
        version: 0.9.0
        maturity: alpha
        replaces: etcdoperator.v0.6.1
        maintainers:
        - name: CoreOS, Inc
          email: support@coreos.com
      
        provider:
          name: CoreOS, Inc
        labels:
          alm-owner-etcd: etcdoperator
          operated-by: etcdoperator
        selector:
          matchLabels:
            alm-owner-etcd: etcdoperator
            operated-by: etcdoperator
        links:
        - name: Blog
          url: https://coreos.com/etcd
        - name: Documentation
          url: https://coreos.com/operators/etcd/docs/latest/
        - name: etcd Operator Source Code
          url: https://github.com/coreos/etcd-operator
      
        icon:
        - base64data: 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
          mediatype: image/png
        install:
          strategy: deployment
          spec:
            permissions:
            - serviceAccountName: etcd-operator
              rules:
              - apiGroups:
                - etcd.database.coreos.com
                resources:
                - etcdclusters
                - etcdbackups
                - etcdrestores
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - pods
                - services
                - endpoints
                - persistentvolumeclaims
                - events
                verbs:
                - "*"
              - apiGroups:
                - apps
                resources:
                - deployments
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - secrets
                verbs:
                - get
            deployments:
            - name: etcd-operator
              spec:
                replicas: 1
                selector:
                  matchLabels:
                    name: etcd-operator-alm-owned
                template:
                  metadata:
                    name: etcd-operator-alm-owned
                    labels:
                      name: etcd-operator-alm-owned
                  spec:
                    serviceAccountName: etcd-operator
                    containers:
                    - name: etcd-operator
                      command:
                      - etcd-operator
                      - --create-crd=false
                      image: quay.io/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
                    - name: etcd-backup-operator
                      image: quay.io/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8
                      command:
                      - etcd-backup-operator
                      - --create-crd=false
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
                    - name: etcd-restore-operator
                      image: quay.io/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8
                      command:
                      - etcd-restore-operator
                      - --create-crd=false
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
        customresourcedefinitions:
          owned:
          - name: etcdclusters.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdCluster
            displayName: etcd Cluster
            description: Represents a cluster of etcd nodes.
            resources:
              - kind: Service
                version: v1
              - kind: Pod
                version: v1
            specDescriptors:
              - description: The desired number of member Pods for the etcd cluster.
                displayName: Size
                path: size
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podCount'
              - description: Limits describes the minimum/maximum amount of compute resources required/allowed
                displayName: Resource Requirements
                path: pod.resources
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:resourceRequirements'
            statusDescriptors:
              - description: The status of each of the member Pods for the etcd cluster.
                displayName: Member Status
                path: members
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podStatuses'
              - description: The service at which the running etcd cluster can be accessed.
                displayName: Service
                path: serviceName
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes:Service'
              - description: The current size of the etcd cluster.
                displayName: Cluster Size
                path: size
              - description: The current version of the etcd cluster.
                displayName: Current Version
                path: currentVersion
              - description: 'The target version of the etcd cluster, after upgrading.'
                displayName: Target Version
                path: targetVersion
              - description: The current status of the etcd cluster.
                displayName: Status
                path: phase
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase'
              - description: Explanation for the current status of the cluster.
                displayName: Status Details
                path: reason
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
          - name: etcdbackups.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdBackup
            displayName: etcd Backup
            description: Represents the intent to backup an etcd cluster.
            specDescriptors:
              - description: Specifies the endpoints of an etcd cluster.
                displayName: etcd Endpoint(s)
                path: etcdEndpoints
                x-descriptors: 
                  - 'urn:alm:descriptor:etcd:endpoint'
              - description: The full AWS S3 path where the backup is saved.
                displayName: S3 Path
                path: s3.path
                x-descriptors: 
                  - 'urn:alm:descriptor:aws:s3:path'
              - description: The name of the secret object that stores the AWS credential and config files.
                displayName: AWS Secret
                path: s3.awsSecret
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:Secret'
            statusDescriptors:
              - description: Indicates if the backup was successful.
                displayName: Succeeded
                path: succeeded
                x-descriptors: 
                  - 'urn:alm:descriptor:text'
              - description: Indicates the reason for any backup related failures.
                displayName: Reason
                path: reason
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
          - name: etcdrestores.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdRestore
            displayName: etcd Restore
            description: Represents the intent to restore an etcd cluster from a backup.
            specDescriptors:
              - description: References the EtcdCluster which should be restored,
                displayName: etcd Cluster
                path: etcdCluster.name
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:EtcdCluster'
                  - 'urn:alm:descriptor:text'
              - description: The full AWS S3 path where the backup is saved.
                displayName: S3 Path
                path: s3.path
                x-descriptors: 
                  - 'urn:alm:descriptor:aws:s3:path'
              - description: The name of the secret object that stores the AWS credential and config files.
                displayName: AWS Secret
                path: s3.awsSecret
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:Secret'
            statusDescriptors:
              - description: Indicates if the restore was successful.
                displayName: Succeeded
                path: succeeded
                x-descriptors: 
                  - 'urn:alm:descriptor:text'
              - description: Indicates the reason for any restore related failures.
                displayName: Reason
                path: reason
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
      
    - #! validate-crd: ./deploy/chart/templates/03-clusterserviceversion.crd.yaml
      #! parse-kind: ClusterServiceVersion
      apiVersion: operators.coreos.com/v1alpha1
      kind: ClusterServiceVersion
      metadata:
        name: etcdoperator.v0.9.2
        namespace: placeholder
        annotations:
          tectonic-visibility: ocs
          alm-examples: '[{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdCluster","metadata":{"name":"example","namespace":"default"},"spec":{"size":3,"version":"3.2.13"}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdRestore","metadata":{"name":"example-etcd-cluster"},"spec":{"etcdCluster":{"name":"example-etcd-cluster"},"backupStorageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}},{"apiVersion":"etcd.database.coreos.com/v1beta2","kind":"EtcdBackup","metadata":{"name":"example-etcd-cluster-backup"},"spec":{"etcdEndpoints":["<etcd-cluster-endpoints>"],"storageType":"S3","s3":{"path":"<full-s3-path>","awsSecret":"<aws-secret>"}}}]'
          packageId: akashem/etcd
          description: Red Hat etcd
          createdAt: 2018/01/01
          certifiedLevel: good
          repository: https://etcd.redhat.com
          containerImage: quay.io/redhat/etcd
          support: Red Hat
          healthIndex: A           
      spec:
        displayName: etcd
        description: |
          etcd is a distributed key value store that provides a reliable way to store data across a cluster of machines. It’s open-source and available on GitHub. etcd gracefully handles leader elections during network partitions and will tolerate machine failure, including the leader. Your applications can read and write data into etcd.
          A simple use-case is to store database connection details or feature flags within etcd as key value pairs. These values can be watched, allowing your app to reconfigure itself when they change. Advanced uses take advantage of the consistency guarantees to implement database leader elections or do distributed locking across a cluster of workers.
      
          _The etcd Open Cloud Service is Public Alpha. The goal before Beta is to fully implement backup features._
      
          ### Reading and writing to etcd
      
          Communicate with etcd though its command line utility etcdctl or with the API using the automatically generated Kubernetes Service.
      
          [Read the complete guide to using the etcd Open Cloud Service](https://coreos.com/tectonic/docs/latest/alm/etcd-ocs.html)
      
          ### Supported Features
      
      
          **High availability**
      
      
          Multiple instances of etcd are networked together and secured. Individual failures or networking issues are transparently handled to keep your cluster up and running.
      
      
          **Automated updates**
      
      
          Rolling out a new etcd version works like all Kubernetes rolling updates. Simply declare the desired version, and the etcd service starts a safe rolling update to the new version automatically.
      
      
          **Backups included**
      
      
          Coming soon, the ability to schedule backups to happen on or off cluster.
        keywords: ['etcd', 'key value', 'database', 'coreos', 'open source']
        version: 0.9.2
        maturity: alpha
        replaces: etcdoperator.v0.9.0
        maintainers:
        - name: CoreOS, Inc
          email: support@coreos.com
      
        provider:
          name: CoreOS, Inc
        labels:
          alm-owner-etcd: etcdoperator
          operated-by: etcdoperator
        selector:
          matchLabels:
            alm-owner-etcd: etcdoperator
            operated-by: etcdoperator
        links:
        - name: Blog
          url: https://coreos.com/etcd
        - name: Documentation
          url: https://coreos.com/operators/etcd/docs/latest/
        - name: etcd Operator Source Code
          url: https://github.com/coreos/etcd-operator
      
        icon:
        - base64data: iVBORw0KGgoAAAANSUhEUgAAAOEAAADZCAYAAADWmle6AAAACXBIWXMAAAsTAAALEwEAmpwYAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAEKlJREFUeNrsndt1GzkShmEev4sTgeiHfRYdgVqbgOgITEVgOgLTEQydwIiKwFQCayoCU6+7DyYjsBiBFyVVz7RkXvqCSxXw/+f04XjGQ6IL+FBVuL769euXgZ7r39f/G9iP0X+u/jWDNZzZdGI/Ftama1jjuV4BwmcNpbAf1Fgu+V/9YRvNAyzT2a59+/GT/3hnn5m16wKWedJrmOCxkYztx9Q+py/+E0GJxtJdReWfz+mxNt+QzS2Mc0AI+HbBBwj9QViKbH5t64DsP2fvmGXUkWU4WgO+Uve2YQzBUGd7r+zH2ZG/tiUQc4QxKwgbwFfVGwwmdLL5wH78aPC/ZBem9jJpCAX3xtcNASSNgJLzUPSQyjB1zQNl8IQJ9MIU4lx2+Jo72ysXYKl1HSzN02BMa/vbZ5xyNJIshJzwf3L0dQhJw4Sih/SFw9Tk8sVeghVPoefaIYCkMZCKbrcP9lnZuk0uPUjGE/KE8JQry7W2tgfuC3vXgvNV+qSQbyFtAtyWk7zWiYevvuUQ9QEQCvJ+5mmu6dTjz1zFHLFj8Eb87MtxaZh/IQFIHom+9vgTWwZxAQjT9X4vtbEVPojwjiV471s00mhAckpwGuCn1HtFtRDaSh6y9zsL+LNBvCG/24ThcxHObdlWc1v+VQJe8LcO0jwtuF8BwnAAUgP9M8JPU2Me+Oh12auPGT6fHuTePE3bLDy+x9pTLnhMn+07TQGh//Bz1iI0c6kvtqInjvPZcYR3KsPVmUsPYt9nFig9SCY8VQNhpPBzn952bbgcsk2EvM89wzh3UEffBbyPqvBUBYQ8ODGPFOLsa7RF096WJ69L+E4EmnpjWu5o4ChlKaRTKT39RMMaVPEQRsz/nIWlDN80chjdJlSd1l0pJCAMVZsniobQVuxceMM9OFoaMd9zqZtjMEYYDW38Drb8Y0DYPLShxn0pvIFuOSxd7YCPet9zk452wsh54FJoeN05hcgSQoG5RR0Qh9Q4E4VvL4wcZq8UACgaRFEQKgSwWrkr5WFnGxiHSutqJGlXjBgIOayhwYBTA0ER0oisIVSUV0AAMT0IASCUO4hRIQSAEECMCCEPwqyQA0JCQBzEGjWNAqHiUVAoXUWbvggOIQCEAOJzxTjoaQ4AIaE64/aZridUsBYUgkhB15oGg1DBIl8IqirYwV6hPSGBSFteMCUBSVXwfYixBmamRubeMyjzMJQBDDowE3OesDD+zwqFoDqiEwXoXJpljB+PvWJGy75BKF1FPxhKygJuqUdYQGlLxNEXkrYyjQ0GbaAwEnUIlLRNvVjQDYUAsJB0HKLE4y0AIpQNgCIhBIhQTgCKhZBBpAN/v6LtQI50JfUgYOnnjmLUFHKhjxbAmdTCaTiBm3ovLPqG2urWAij6im0Nd9aTN9ygLUEt9LgSRnohxUPIKxlGaE+/6Y7znFf0yX+GnkvFFWmarkab2o9PmTeq8sbd2a7DaysXz7i64VeznN4jCQhN9gdDbRiuWrfrsq0mHIrlaq+hlotCtd3Um9u0BYWY8y5D67wccJoZjFca7iUs9VqZcfsZwTd1sbWGG+OcYaTnPAP7rTQVVlM4Sg3oGvB1tmNh0t/HKXZ1jFoIMwCQjtqbhNxUmkGYqgZEDZP11HN/S3gAYRozf0l8C5kKEKUvW0t1IfeWG/5MwgheZTT1E0AEhDkAePQO+Ig2H3DncAkQM4cwUQCD530dU4B5Yvmi2LlDqXfWrxMCcMth51RToRMNUXFnfc2KJ0+Ryl0VNOUwlhh6NoxK5gnViTgQpUG4SqSyt5z3zRJpuKmt3Q1614QaCBPaN6je+2XiFcWAKOXcUfIYKRyL/1lb7pe5VxSxxjQ6hImshqGRt5GWZVKO6q2wHwujfwDtIvaIdexj8Cm8+a68EqMfox6x/voMouZF4dHnEGNeCDMwT6vdNfekH1MafMk4PI06YtqLVGl95aEM9Z5vAeCTOA++YLtoVJRrsqNCaJ6WRmkdYaNec5BT/lcTRMqrhmwfjbpkj55+OKp8IEbU/JLgPJE6Wa3TTe9sHS+ShVD5QIyqIxMEwKh12olC6mHIed5ewEop80CNlfIOADYOT2nd6ZXCop+Ebqchc0JqxKcKASxChycJgUh1rnHA5ow9eTrhqNI7JWiAYYwBGGdpyNLoGw0Pkh96h1BpHihyywtATDM/7Hk2fN9EnH8BgKJCU4ooBkbXFMZJiPbrOyecGl3zgQDQL4hk10IZiOe+5w99Q/gBAEIJgPhJM4QAEEoFREAIAAEiIASAkD8Qt4AQAEIAERAGFlX4CACKAXGVM4ivMwWwCLFAlyeoaa70QePKm5Dlp+/n+ye/5dYgva6YsUaVeMa+tzNFeJtWwc+udbJ0Fg399kLielQJ5Ze61c2+7ytA6EZetiPxZC6tj22yJCv6jUwOyj/zcbqAxOMyAKEbfeHtNa7DtYXptjsk2kJxR+eIeim/tHNofUKYy8DMrQcAKWz6brpvzyIAlpwPhQ49l6b7skJf5Z+YTOYQc4FwLDxvoTDwaygQK+U/kVr+ytSFBG01Q3gnJJR4cNiAhx4HDub8/b5DULXlj6SVZghFiE+LdvE9vo/o8Lp1RmH5hzm0T6wdbZ6n+D6i44zDRc3ln6CpAEJfXiRU45oqLz8gFAThWsh7ughrRibc0QynHgZpNJa/ENJ+loCwu/qOGnFIjYR/n7TfgycULhcQhu6VC+HfF+L3BoAQ4WiZTw1M+FPCnA2gKC6/FAhXgDC+ojQGh3NuWsvfF1L/D5ohlCKtl1j2ldu9a/nPAKFwN56Bst10zCG0CPleXN/zXPgHQZXaZaBgrbzyY5V/mUA+6F0hwtGN9rwu5DVZPuwWqfxdFz1LWbJ2lwKEa+0Qsm4Dl3fp+Pu0lV97PgwIPfSsS+UQhj5Oo+vvFULazRIQyvGEcxPuNLCth2MvFsrKn8UOilAQShkh7TTczYNMoS6OdP47msrPi82lXKGWhCdMZYS0bFy+vcnGAjP1CIfvgbKNA9glecEH9RD6Ol4wRuWyN/G9MHnksS6o/GPf5XcwNSUlHzQhDuAKtWJmkwKElU7lylP5rgIcsquh/FI8YZCDpkJBuE4FQm7Icw8N+SrUGaQKyi8FwiDt1ve5o+Vu7qYHy/psgK8cvh+FTYuO77bhEC7GuaPiys/L1X4IgXDL+e3M5+ovLxBy5VLuIebw1oqcHoPfoaMJUsHays878r8KbDc3xtPx/84gZPBG/JwaufrsY/SRG/OY3//8QMNdsvdZCFtbW6f8pFuf5bflILAlX7O+4fdfugKyFYS8T2zAsXthdG0VurPGKwI06oF5vkBgHWkNp6ry29+lsPZMU3vijnXFNmoclr+6+Ou/FIb8yb30sS8YGjmTqCLyQsi5N/6ZwKs0Yenj68pfPjF6N782Dp2FzV9CTyoSeY8mLK16qGxIkLI8oa1n8tz9juP40DlK0epxYEbojbq+9QfurBeVIlCO9D2396bxiV4lkYQ3hOAFw2pbhqMGISkkQOMcQ9EqhDmGZZdo92JC0YHRNTfoSg+5e0IT+opqCKHoIU+4ztQIgBD1EFNrQAgIpYSil9lDmPHqkROPt+JC6AgPquSuumJmg0YARVCuneDfvPVeJokZ6pIXDkNxQtGzTF9/BQjRG0tQznfb74RwCQghpALBtIQnfK4zhxdyQvVCUeknMIT3hLyY+T5jo0yABqKPQNpUNw/09tGZod5jgCaYFxyYvJcNPkv9eof+I3pnCFEHIETjSM8L9tHZHYCQT9PaZGycU6yg8S4akDnJ+P03L0+t23XGzCLzRgII/Wqa+fv/xlfvmKvMUOcOrlCDdoei1MGdZm6G5VEIfRzzjd4aQs69n699Rx7ewhvCGzr2gmTPs8zNsJOrXt24FbkhhOjCfT4ICA/rPbyhUy94Dks0gJCX1NzCZui9YUd3oei+c257TalFbgg19ILHrlrL2gvWgXAL26EX76gZTNASQnad8Ibwhl284NhgXpB0c+jKhWO3Ms1hP9ihJYB9eMF6qd1BCPk0qA1s+LimFIu7m4nsdQIzPK4VbQ8hYvrnuSH2G9b2ggP78QmWqBdF9Vx8SSY6QYdUW7BTA1schZATyhvY8lHvcRbNUS9YGFy2U+qmzh2YPVc0I7yAOFyHfRpyUwtCSzOdPXMHmz7qDIM0e0V2wZTEk+6Ym6N63eBLp/b5Bts+2cKCSJ/LuoZO3ANSiE5hKAZjnvNSS4931jcw9jpwT0feV/qSJ1pVtCyfHKDkvK8Ejx7pUxGh2xFNSwx8QTi2H9ceC0/nni64MS/5N5dG39pDqvRV+WgGk71c9VFXF9b+xYvOw/d61iv7m3MvEHryhvecwC52jSSx4VIIgwnMNT/UsTxIgpPt3K/ARj15CptwL3Zd/ceDSATj2DGQjbxgWwhdeMMte7zpy5On9vymRm/YxBYljGVjKWF9VJf7I1+sex3wY8w/V1QPTborW/72gkdsRDaZMJBdbdHIC7aCkAu9atlLbtnrzerMnyToDaGwelOnk3/hHSem/ZK7e/t7jeeR20LYBgqa8J80gS8jbwi5F02Uj1u2NYJxap8PLkJfLxA2hIJyvnHX/AfeEPLpBfe0uSFHbnXaea3Qd5d6HcpYZ8L6M7lnFwMQ3MNg+RxUR1+6AshtbsVgfXTEg1sIGax9UND2p7f270wdG3eK9gXVGHdw2k5sOyZv+Nbs39Z308XR9DqWb2J+PwKDhuKHPobfuXf7gnYGHdCs7bhDDadD4entDug7LWNsnRNW4mYqwJ9dk+GGSTPBiA2j0G8RWNM5upZtcG4/3vMfP7KnbK2egx6CCnDPhRn7NgD3cghLIad5WcM2SO38iqHvvMOosyeMpQ5zlVCaaj06GVs9xUbHdiKoqrHWgquFEFMWUEWfXUxJAML23hAHFOctmjZQffKD2pywkhtSGHKNtpitLroscAeE7kCkSsC60vxEl6yMtL9EL5HKGCMszU5bk8gdkklAyEn5FO0yK419rIxBOIqwFMooDE0tHEVYijAUECIshRCGIhxFWIowFJ5QkEYIS5PTJrUwNGlPyN6QQPyKtpuM1E/K5+YJDV/MiA3AaehzqgAm7QnZG9IGYKo8bHnSK7VblLL3hOwNHziPuEGOqE5brrdR6i+atCfckyeWD47HkAkepRGLY/e8A8J0gCwYSNypF08bBm+e6zVz2UL4AshhBUjML/rXLefqC82bcQFhGC9JDwZ1uuu+At0S5gCETYHsV4DUeD9fDN2Zfy5OXaW2zAwQygCzBLJ8cvaW5OXKC1FxfTggFAHmoAJnSiOw2wps9KwRWgJCLaEswaj5NqkLwAYIU4BxqTSXbHXpJdRMPZgAOiAMqABCNGYIEEJutEK5IUAIwYMDQgiCACEEAcJs1Vda7gGqDhCmoiEghAAhBAHCrKXVo2C1DCBMRlp37uMIEECoX7xrX3P5C9QiINSuIcoPAUI0YkAICLNWgfJDh4T9hH7zqYH9+JHAq7zBqWjwhPAicTVCVQJCNF50JghHocahKK0X/ZnQKyEkhSdUpzG8OgQI42qC94EQjsYLRSmH+pbgq73L6bYkeEJ4DYTYmeg1TOBFc/usTTp3V9DdEuXJ2xDCUbXhaXk0/kAYmBvuMB4qkC35E5e5AMKkwSQgyxufyuPy6fMMgAFCSI73LFXU/N8AmEL9X4ABACNSKMHAgb34AAAAAElFTkSuQmCC
          mediatype: image/png
        install:
          strategy: deployment
          spec:
            permissions:
            - serviceAccountName: etcd-operator
              rules:
              - apiGroups:
                - etcd.database.coreos.com
                resources:
                - etcdclusters
                - etcdbackups
                - etcdrestores
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - pods
                - services
                - endpoints
                - persistentvolumeclaims
                - events
                verbs:
                - "*"
              - apiGroups:
                - apps
                resources:
                - deployments
                verbs:
                - "*"
              - apiGroups:
                - ""
                resources:
                - secrets
                verbs:
                - get
            deployments:
            - name: etcd-operator
              spec:
                replicas: 1
                selector:
                  matchLabels:
                    name: etcd-operator-alm-owned
                template:
                  metadata:
                    name: etcd-operator-alm-owned
                    labels:
                      name: etcd-operator-alm-owned
                  spec:
                    serviceAccountName: etcd-operator
                    containers:
                    - name: etcd-operator
                      command:
                      - etcd-operator
                      - --create-crd=false
                      image: quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
                    - name: etcd-backup-operator
                      image: quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2
                      command:
                      - etcd-backup-operator
                      - --create-crd=false
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
                    - name: etcd-restore-operator
                      image: quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2
                      command:
                      - etcd-restore-operator
                      - --create-crd=false
                      env:
                      - name: MY_POD_NAMESPACE
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.namespace
                      - name: MY_POD_NAME
                        valueFrom:
                          fieldRef:
                            fieldPath: metadata.name
        customresourcedefinitions:
          owned:
          - name: etcdclusters.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdCluster
            displayName: etcd Cluster
            description: Represents a cluster of etcd nodes.
            resources:
              - kind: Service
                version: v1
              - kind: Pod
                version: v1
            specDescriptors:
              - description: The desired number of member Pods for the etcd cluster.
                displayName: Size
                path: size
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podCount'
              - description: Limits describes the minimum/maximum amount of compute resources required/allowed
                displayName: Resource Requirements
                path: pod.resources
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:resourceRequirements'
            statusDescriptors:
              - description: The status of each of the member Pods for the etcd cluster.
                displayName: Member Status
                path: members
                x-descriptors:
                  - 'urn:alm:descriptor:com.tectonic.ui:podStatuses'
              - description: The service at which the running etcd cluster can be accessed.
                displayName: Service
                path: serviceName
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes:Service'
              - description: The current size of the etcd cluster.
                displayName: Cluster Size
                path: size
              - description: The current version of the etcd cluster.
                displayName: Current Version
                path: currentVersion
              - description: 'The target version of the etcd cluster, after upgrading.'
                displayName: Target Version
                path: targetVersion
              - description: The current status of the etcd cluster.
                displayName: Status
                path: phase
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase'
              - description: Explanation for the current status of the cluster.
                displayName: Status Details
                path: reason
                x-descriptors:
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
          - name: etcdbackups.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdBackup
            displayName: etcd Backup
            description: Represents the intent to backup an etcd cluster.
            specDescriptors:
              - description: Specifies the endpoints of an etcd cluster.
                displayName: etcd Endpoint(s)
                path: etcdEndpoints
                x-descriptors: 
                  - 'urn:alm:descriptor:etcd:endpoint'
              - description: The full AWS S3 path where the backup is saved.
                displayName: S3 Path
                path: s3.path
                x-descriptors: 
                  - 'urn:alm:descriptor:aws:s3:path'
              - description: The name of the secret object that stores the AWS credential and config files.
                displayName: AWS Secret
                path: s3.awsSecret
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:Secret'
            statusDescriptors:
              - description: Indicates if the backup was successful.
                displayName: Succeeded
                path: succeeded
                x-descriptors: 
                  - 'urn:alm:descriptor:text'
              - description: Indicates the reason for any backup related failures.
                displayName: Reason
                path: reason
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
          - name: etcdrestores.etcd.database.coreos.com
            version: v1beta2
            kind: EtcdRestore
            displayName: etcd Restore
            description: Represents the intent to restore an etcd cluster from a backup.
            specDescriptors:
              - description: References the EtcdCluster which should be restored,
                displayName: etcd Cluster
                path: etcdCluster.name
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:EtcdCluster'
                  - 'urn:alm:descriptor:text'
              - description: The full AWS S3 path where the backup is saved.
                displayName: S3 Path
                path: s3.path
                x-descriptors: 
                  - 'urn:alm:descriptor:aws:s3:path'
              - description: The name of the secret object that stores the AWS credential and config files.
                displayName: AWS Secret
                path: s3.awsSecret
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes:Secret'
            statusDescriptors:
              - description: Indicates if the restore was successful.
                displayName: Succeeded
                path: succeeded
                x-descriptors: 
                  - 'urn:alm:descriptor:text'
              - description: Indicates the reason for any restore related failures.
                displayName: Reason
                path: reason
                x-descriptors: 
                  - 'urn:alm:descriptor:io.kubernetes.phase:reason'
  packages: |-      
    - #! package-manifest: ./deploy/chart/catalog_resources/rh-operators/etcdoperator.v0.9.2.clusterserviceversion.yaml
      packageName: etcd
      channels:
      - name: alpha
        currentCSV: etcdoperator.v0.9.2
`
