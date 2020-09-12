
class B

class A
  field B

entity:
  id as an identifier

value object:
  value as an identifier

A is depend on B means: 
  A is defined by B value
  B may independent (can be used by other) or only defined under A

A is not depend on B means: 
  A is fully defined even without B value. 
  B is only an additional information

B is depend on A, means: 
  B does not exist if there is no A

B is not depend on A, means: 
  B is independent

if A is not depend on B and B is not depend on A
  A have ref to B (mandatory, since B is inside A)

if A is not depend on B and B is depend on A
  B have ref to A (mandatory)
  A have ref to B (optional)

if A depend on B and B is not depend on A
  A have ref to B (mandatory)

if A depend on B and B is depend on A
  A have ref to B (mandatory)
  B have ref to A (mandatory)

==================================

interface C

class B : C

class D : C

class A
  field C

