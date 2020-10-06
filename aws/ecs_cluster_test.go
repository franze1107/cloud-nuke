package aws

import (
	"fmt"
	"testing"
	"time"

	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/cloud-nuke/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const tagKey = "first_seen"

// Test that we can succesfully list ECS clusters by manually creating a cluster and then using the list function to find it.
func TestCanCreateAndListEcsCluster(t *testing.T) {
	t.Parallel()

	region := "eu-west-1"
	awsSession, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region),
	})
	require.NoError(t, err)

	clusterName := fmt.Sprintf("test-ecs-cluster-%s", util.UniqueID())
	cluster := createEcsFargateCluster(t, awsSession, clusterName)
	defer deleteEcsCluster(awsSession, cluster)

	clusterArns, err := getAllEcsClusters(awsSession)
	require.NoError(t, err)

	assert.Contains(t, clusterArns, cluster.ClusterArn)
}

// Test we can create a cluster, tag it, and then find the tag
func TestCanTagEcsClusters(t *testing.T) {
	t.Parallel()

	region := "eu-west-1"
	awsSession, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region),
	})
	require.NoError(t, err)

	cluster := createEcsFargateCluster(t, awsSession, util.UniqueID())
	defer deleteEcsCluster(awsSession, cluster)

	tagValue := time.Now().UTC().Format(time.RFC3339)

	tagEcsCluster(awsSession, cluster.ClusterArn, "first_seen", tagValue)
	require.NoError(t, err)

	returnedTag, err := getClusterTag(awsSession, cluster.ClusterArn, "first_seen")
	require.NoError(t, err)

	assert.Equal(t, returnedTag.Format(time.RFC3339), tagValue)
}

// Test that we can filter ECS clusters by 'created_at' tag value.
func TestCanFilterOlderEcsClusters(t *testing.T) {
	t.Parallel()

	region := "eu-west-1"
	awsSession, err := session.NewSession(&awsgo.Config{
		Region: awsgo.String(region),
	})
	require.NoError(t, err)

	clusterName := util.UniqueID()
	cluster := createEcsFargateCluster(t, awsSession, clusterName)
	defer deleteEcsCluster(awsSession, cluster)

	tagEcsCluster(awsSession, cluster.ClusterArn, tagKey, time.Now().UTC().String())
	require.NoError(t, err)

}

// Test we can get all ECS clusters younger than < X time based on tags
func TestCanListAllEcsClustersOlderThan24hours(t *testing.T) {
	// create 3 clusters with tags: 1hr, 22hrs, 28hrs
	// get all ecs clusters
	// get tags for each cluster
	// select only clusters older than 24hrs
	// assert return only 1 cluster
}

// Test we can nuke all ECS clusters younger than < X time
func TestCanNukeAllEcsClustersOlderThan24Hours(t *testing.T) {
	// create 3 clusters with tags: 1hr, 25hrs, 28hrs
	// get all ecs clusters
	// get tags for each cluster
	// select only clusters older than 24hrs
	// nuke selected clusters
	// assert 2 clusters nuked
	// assert 1 cluster still left
}
