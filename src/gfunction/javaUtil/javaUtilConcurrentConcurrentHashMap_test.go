package javaUtil

import (
	"jacobin/src/gfunction/ghelpers"
	"reflect"
	"testing"
)

func TestLoad_Util_Concurrent_Hash_Map(t *testing.T) {
	saved := ghelpers.MethodSignatures
	defer func() { ghelpers.MethodSignatures = saved }()
	ghelpers.MethodSignatures = make(map[string]ghelpers.GMeth)

	Load_Util_Concurrent_Hash_Map()

	expectedMethods := map[string]struct {
		slots int
		trap  any
	}{
		"java/util/concurrent/ConcurrentHashMap.<clinit>()V":                                                                                              {0, ghelpers.TrapClass},
		"java/util/concurrent/ConcurrentHashMap.<init>()V":                                                                                                {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.<init>(I)V":                                                                                               {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.<init>(IF)V":                                                                                              {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.<init>(IFI)V":                                                                                             {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.<init>(Ljava/util/Map;)V":                                                                                 {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.newKeySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                           {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.newKeySet(I)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                          {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.clear()V":                                                                                                 {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.compute(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                             {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.computeIfAbsent(Ljava/lang/Object;Ljava/util/function/Function;)Ljava/lang/Object;":                       {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.computeIfPresent(Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                    {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.contains(Ljava/lang/Object;)Z":                                                                            {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.containsKey(Ljava/lang/Object;)Z":                                                                         {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.containsValue(Ljava/lang/Object;)Z":                                                                       {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.elements()Ljava/util/Enumeration;":                                                                        {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.entrySet()Ljava/util/Set;":                                                                                {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.equals(Ljava/lang/Object;)Z":                                                                              {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEach(Ljava/util/function/BiConsumer;)V":                                                                {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiConsumer;)V":                                                               {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEach(JLjava/util/function/BiFunction;Ljava/util/function/Consumer;)V":                                  {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Consumer;)V":                                                            {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachEntry(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                               {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Consumer;)V":                                                              {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachKey(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                                 {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Consumer;)V":                                                            {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.forEachValue(JLjava/util/function/Function;Ljava/util/function/Consumer;)V":                               {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.get(Ljava/lang/Object;)Ljava/lang/Object;":                                                                {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.getOrDefault(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                     {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.hashCode()I":                                                                                              {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.isEmpty()Z":                                                                                               {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/Set;":                                                                                  {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.keySet()Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                                              {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.keySet(Ljava/lang/Object;)Ljava/util/concurrent/ConcurrentHashMap$KeySetView;":                            {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.keys()Ljava/util/Enumeration;":                                                                            {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.mappingCount()J":                                                                                          {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.merge(Ljava/lang/Object;Ljava/lang/Object;Ljava/util/function/BiFunction;)Ljava/lang/Object;":             {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.put(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                              {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.putAll(Ljava/util/Map;)V":                                                                                 {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.putIfAbsent(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                      {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduce(JLjava/util/function/BiFunction;Ljava/util/function/BiFunction;)Ljava/lang/Object;":                {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/BiFunction;)Ljava/util/Map$Entry;":                                     {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntries(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":           {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D": {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":          {5, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceEntriesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":       {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                           {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeys(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":              {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D":    {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":             {5, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceKeysToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":          {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToDouble(JLjava/util/function/ToDoubleBiFunction;DLjava/util/function/DoubleBinaryOperator;)D":      {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToInt(JLjava/util/function/ToIntBiFunction;ILjava/util/function/IntBinaryOperator;)I":               {5, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceToLong(JLjava/util/function/ToLongBiFunction;JLjava/util/function/LongBinaryOperator;)J":            {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                         {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValues(JLjava/util/function/Function;Ljava/util/function/BiFunction;)Ljava/lang/Object;":            {4, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToDouble(JLjava/util/function/ToDoubleFunction;DLjava/util/function/DoubleBinaryOperator;)D":  {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToInt(JLjava/util/function/ToIntFunction;ILjava/util/function/IntBinaryOperator;)I":           {5, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.reduceValuesToLong(JLjava/util/function/ToLongFunction;JLjava/util/function/LongBinaryOperator;)J":        {6, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;)Ljava/lang/Object;":                                                             {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.remove(Ljava/lang/Object;Ljava/lang/Object;)Z":                                                            {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;":                                          {2, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.replace(Ljava/lang/Object;Ljava/lang/Object;Ljava/lang/Object;)Z":                                         {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.replaceAll(Ljava/util/function/BiFunction;)V":                                                             {1, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.search(JLjava/util/function/BiFunction;)Ljava/lang/Object;":                                               {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchEntries(JLjava/util/function/Function;)Ljava/lang/Object;":                                          {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchKeys(JLjava/util/function/Function;)Ljava/lang/Object;":                                             {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.searchValues(JLjava/util/function/Function;)Ljava/lang/Object;":                                           {3, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.size()I":                                                                                                  {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.toString()Ljava/lang/String;":                                                                             {0, ghelpers.TrapFunction},
		"java/util/concurrent/ConcurrentHashMap.values()Ljava/util/Collection;":                                                                           {0, ghelpers.TrapFunction},
	}

	if len(ghelpers.MethodSignatures) != len(expectedMethods) {
		t.Errorf("Expected %d registered methods, got %d", len(expectedMethods), len(ghelpers.MethodSignatures))
	}

	for sig, exp := range expectedMethods {
		gm, ok := ghelpers.MethodSignatures[sig]
		if !ok {
			t.Errorf("Missing registration for %s", sig)
			continue
		}
		if gm.ParamSlots != exp.slots {
			t.Errorf("%s: ParamSlots expected %d, got %d", sig, exp.slots, gm.ParamSlots)
		}
		if gm.GFunction == nil {
			t.Errorf("%s: GFunction is nil", sig)
		}
		if reflect.ValueOf(gm.GFunction).Pointer() != reflect.ValueOf(exp.trap).Pointer() {
			t.Errorf("%s: GFunction mismatch", sig)
		}
	}
}
