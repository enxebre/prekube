package main

type Cloud interface {
	elected() bool
	retrieveAssets()
}

func main() {
	cloudImpl := newGcp(
		"breaking-170711",
		"europe-west1",
		"kube-igm",
		"kube",
		"spiderman.png",
		"/Users/albertogarla/Desktop/content.png",
	)
	preKube(cloudImpl)
}

func preKube(c Cloud) {
	if c.elected() {
		c.retrieveAssets()
	}
	return
}
