

def sorted(L):
	if not L: return L
	
  	L1 = [i for i in L[1:] if i < L[0]]
	L2 = [i for i in L[1:] if i >= L[0]]
	return sorted(L1) + [L[0]] + sorted(L2)

print sorted([1,3,2,4,5,6,3,2,1,4,6])