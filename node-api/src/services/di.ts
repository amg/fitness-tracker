class DI {
    constructor() {}
    
    private registry = new Map<string, Dependency<any>>();

    register<T>(key: string, factory: () => T) {
        this.registry.get(key) != null ? console.log(`[WARN] DI.register: redeclared dependency: ${key}`): void
        this.registry.set(key, new Dependency<T>(null, factory))
    }

    get<T>(key:string): T | null {
        const dep = this.registry.get(key);
        if (dep == null) {
            console.log(`[WARN] DI.get: dependency is not declared`);
            return null
        }
        if (dep.value == null) {
            dep.value = dep!!.factory()
        }
        return dep.value
    }
}

class Dependency<T> {
    constructor(public value: T | null, public factory: () => T) {}
}