# Inovkedynamic Status

jvm.interpret() -> jvm.doInvokeDynamic() -> 
mhResoution.go::jvm.resolveCallSite() does:
1. Fetch the InvokeDynamic entry -- OK
2. Get the Bootstrap Method info from the class attributes -- OK
3. Resolve the Bootstrap Method Handle -> classloader.ResolveMethodHandle() 
