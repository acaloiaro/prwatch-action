package internal

// serviceProviders is a functional seam that enables service provider implementations to be easily swapped out
// application-wide.
// serviceProviders should only reference interfaces; not implementations.
type serviceProvider struct {
	g gitProvider
	f fileProvider
	i issueProvider
}

// the global services provider for all prwatch
var services = newProvider()

func newProvider() serviceProvider {
	return serviceProvider{}
}

// TODO: reset() is useful for resetting providers during testing. However, I'm not fond of having test helper code adjacent
// to the implementation.
func (p serviceProvider) reset() {
	services = newProvider()
}

func (p serviceProvider) git() gitProvider {
	if p.g == nil {
		p.g = &GitCommandLine{}
	}

	return p.g
}

func (p serviceProvider) files() fileProvider {
	if p.f == nil {
		p.f = &posixFileProvider{}
	}

	return p.f
}

func (p serviceProvider) issues() issueProvider {
	if p.i == nil {
		p.i = newJiraIssueProvider(newJiraClient())
	}

	return p.i
}
