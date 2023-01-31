# Ginkgo & Gomega Demo

### Installation
```
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Widely used flags

```
ginkgo
ginkgo ./...
ginkgo -r
ginkgo -p -procs=N
ginkgo --race
ginkgo --randomize-all --randomize-suites
ginkgo --until-it-fails
ginkgo --repeat=N
ginkgo --timeout=duration //default 1h
ginkgo --skip-package=list,of,packages
ginkgo --coverprofile=coverage.out

```

### Widely used matchers
Matching objects:

``` 
Expect(ACTUAL).To(Equal(EXPECTED))
Expect(ACTUAL).ToNot(Equal(EXPECTED))
```

Checking error
``` 
Expect(ERROR).To(HaveOccurred())
Expect(ERROR).ToNot(HaveOccurred())
```    

Checking object type - good use case is checking error type
``` 
Expect(ACTUAL).To(BeAssignableToTypeOf(TYPE))
```    

Checking for empty object - string, array, map, slice, channel
``` 
Expect(ERROR).Should(BeEmpty())
```    

Checking boolean
``` 
Expect(BOOL).Should(BeTrue())
Expect(BOOL).Should(BeFalse())
```   

Expect a behavior to occur within a time
``` 
// if timeout or polling interval is not provided, a default one will be used
Eventually(ACTUAL, (TIMEOUT), (POLLING_INTERVAL), (context.Context)).Should(Equal(EXPECTED))
```   

Expect a behavior to occur for a given time
``` 
Consistently(ACTUAL, (TIMEOUT), (POLLING_INTERVAL), (context.Context)).Should(Equal(EXPECTED))
```   



### Useful links
- Ginkgo documentation: https://onsi.github.io/ginkgo
- Ginkgo github: https://github.com/onsi/ginkgo
- Gomega documentation: https://onsi.github.io/gomega