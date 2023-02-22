## Interface as Contract

In the context of object-oriented programming, an interface is a way of defining a contract between two or more components in a system. An interface specifies a set of methods, properties, and events that a class must implement, but it doesn't provide an implementation of those methods. Instead, it simply defines the method signatures and their return types, and leaves the actual implementation of those methods up to the class that implements the interface.

The idea of an interface as a contract is that it defines what a class is supposed to do, but not how it is supposed to do it. The contract specifies what inputs and outputs are expected, and what behavior is required of the class, but leaves the details of how the class should work up to the individual class. This allows for greater flexibility and modularity in a system, as different components can be swapped in and out as long as they adhere to the same interface.


## Fat Interface

In programming, a "fat interface" refers to an interface that has too many methods, which are not necessarily related to each other. A fat interface violates the Single Responsibility Principle (SRP) because it combines multiple responsibilities into a single interface.
When an interface is overloaded with too many methods or parameters, it can become unclear what the primary purpose of the interface is, and it may be difficult for developers to use or extend the component in a consistent way.

## Interface Segregation Principle

To solve the problem of a fat interface, the Interface Segregation Principle (ISP) can be applied. The Interface Segregation Principle states that a client should not be forced to depend on methods that it does not use. In other words, an interface should be divided into smaller, more specific interfaces, so that each client only needs to depend on the methods it actually uses.

