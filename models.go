package main

type Aws struct {
	Items []struct {
		Clusters struct {
			L []struct {
				M struct {
					ClusterVersion struct {
						S string `json:"S"`
					} `json:"cluster_version"`
					Active struct {
						Bool bool `json:"BOOL"`
						S    bool `json:"S"`
					} `json:"active"`
					EksVersion struct {
						S string `json:"S"`
					} `json:"eks_version"`
					EnableThanos struct {
						S bool `json:"S"`
					} `json:"enable_thanos"`
					DeploymentDate struct {
						S string `json:"S"`
					} `json:"deployment_date"`
					AmiNameRegex struct {
						S string `json:"S"`
					} `json:"ami_name_regex"`
				} `json:"M"`
			} `json:"L"`
		} `json:"clusters"`
	} `json:"Items"`
	Count            int         `json:"Count"`
	ScannedCount     int         `json:"ScannedCount"`
	ConsumedCapacity interface{} `json:"ConsumedCapacity"`
}
