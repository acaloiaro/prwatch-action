package internal

// serviceProviders is a functional seam that enables interface implementations to be easily swapped out
// application-wide.
// serviceProviders should only reference interfaces; not implementations.
type serviceProvider struct {
	git   gitProvider
	files fileProvider
}

var services = &serviceProvider{
	git:   &GitCommandLine{},
	files: &posixFileProvider{},
}

func (p serviceProvider) reset() {
	services = &serviceProvider{
		git:   &GitCommandLine{},
		files: &posixFileProvider{},
	}
}
