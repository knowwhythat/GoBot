export namespace forms {
	
	export class ExecutionForm {
	    id: string;
	    projectId: string;
	    projectName: string;
	    subFlowId: string;
	    triggerType: string;
	    executeResult: number;
	    errorMsg: string;
	    // Go type: time
	    startTs: any;
	    // Go type: time
	    endTs: any;
	
	    static createFrom(source: any = {}) {
	        return new ExecutionForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.projectName = source["projectName"];
	        this.subFlowId = source["subFlowId"];
	        this.triggerType = source["triggerType"];
	        this.executeResult = source["executeResult"];
	        this.errorMsg = source["errorMsg"];
	        this.startTs = this.convertValues(source["startTs"], null);
	        this.endTs = this.convertValues(source["endTs"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ImageData {
	    id: string;
	    name: string;
	    image: string;
	    left: number;
	    top: number;
	    right: number;
	    bottom: number;
	
	    static createFrom(source: any = {}) {
	        return new ImageData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.image = source["image"];
	        this.left = source["left"];
	        this.top = source["top"];
	        this.right = source["right"];
	        this.bottom = source["bottom"];
	    }
	}
	export class LoginForm {
	    username: string;
	    pwd: string;
	    rememberMe: boolean;
	    autoLogin: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LoginForm(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.pwd = source["pwd"];
	        this.rememberMe = source["rememberMe"];
	        this.autoLogin = source["autoLogin"];
	    }
	}
	export class Package {
	    name: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new Package(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	    }
	}
	export class PipPackage {
	    name: string;
	    upgradePip: boolean;
	    upgrade: boolean;
	    inputVersion: boolean;
	    version: string;
	    inputMirror: boolean;
	    mirrorUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new PipPackage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.upgradePip = source["upgradePip"];
	        this.upgrade = source["upgrade"];
	        this.inputVersion = source["inputVersion"];
	        this.version = source["version"];
	        this.inputMirror = source["inputMirror"];
	        this.mirrorUrl = source["mirrorUrl"];
	    }
	}
	export class RunningInstance {
	    id: string;
	    projectId: string;
	    projectName: string;
	    // Go type: time
	    startTs: any;
	    triggerType: string;
	    isVirtual: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RunningInstance(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.projectName = source["projectName"];
	        this.startTs = this.convertValues(source["startTs"], null);
	        this.triggerType = source["triggerType"];
	        this.isVirtual = source["isVirtual"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace models {
	
	export class Project {
	    id: string;
	    name: string;
	    path: string;
	    icon: string;
	    isFlow: boolean;
	    inputParam: {[key: string]: any};
	    description: string;
	    // Go type: time
	    createTs: any;
	    // Go type: time
	    updateTs: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.icon = source["icon"];
	        this.isFlow = source["isFlow"];
	        this.inputParam = source["inputParam"];
	        this.description = source["description"];
	        this.createTs = this.convertValues(source["createTs"], null);
	        this.updateTs = this.convertValues(source["updateTs"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Variable {
	    name: string;
	    type: string;
	    value: string;
	    kind: string;
	
	    static createFrom(source: any = {}) {
	        return new Variable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.value = source["value"];
	        this.kind = source["kind"];
	    }
	}
	export class Flow {
	    key: string;
	    label: string;
	    nodeType: string;
	    opened: boolean;
	    children: Flow[];
	
	    static createFrom(source: any = {}) {
	        return new Flow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.nodeType = source["nodeType"];
	        this.opened = source["opened"];
	        this.children = this.convertValues(source["children"], Flow);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ProjectConfig {
	    key: string;
	    label: string;
	    nodeType: string;
	    opened: boolean;
	    children: Flow[];
	    variables: Variable[];
	    externalDependencies: string[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.nodeType = source["nodeType"];
	        this.opened = source["opened"];
	        this.children = this.convertValues(source["children"], Flow);
	        this.variables = this.convertValues(source["variables"], Variable);
	        this.externalDependencies = source["externalDependencies"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Schedule {
	    id: string;
	    projectId: string;
	    projectVersionId: string;
	    name: string;
	    cron: string;
	    enabled: boolean;
	    desc: string;
	    config: string;
	    // Go type: time
	    createTs: any;
	    // Go type: time
	    updateTs: any;
	
	    static createFrom(source: any = {}) {
	        return new Schedule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.projectId = source["projectId"];
	        this.projectVersionId = source["projectVersionId"];
	        this.name = source["name"];
	        this.cron = source["cron"];
	        this.enabled = source["enabled"];
	        this.desc = source["desc"];
	        this.config = source["config"];
	        this.createTs = this.convertValues(source["createTs"], null);
	        this.updateTs = this.convertValues(source["updateTs"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace plugin {
	
	export class Output {
	    key: string;
	    label: string;
	    placeholder: string;
	    type: string;
	    default_value: string;
	    required: boolean;
	    editor_type: string;
	
	    static createFrom(source: any = {}) {
	        return new Output(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.placeholder = source["placeholder"];
	        this.type = source["type"];
	        this.default_value = source["default_value"];
	        this.required = source["required"];
	        this.editor_type = source["editor_type"];
	    }
	}
	export class Input {
	    key: string;
	    label: string;
	    placeholder: string;
	    category?: string;
	    type: string;
	    default_value: string;
	    required: boolean;
	    options: any;
	    editor_type: string;
	    show_if?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Input(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.placeholder = source["placeholder"];
	        this.category = source["category"];
	        this.type = source["type"];
	        this.default_value = source["default_value"];
	        this.required = source["required"];
	        this.options = source["options"];
	        this.editor_type = source["editor_type"];
	        this.show_if = source["show_if"];
	    }
	}
	export class ParameterDefine {
	    error?: boolean;
	    inputs?: Input[];
	    extra?: Input[];
	    outputs?: Output[];
	
	    static createFrom(source: any = {}) {
	        return new ParameterDefine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.error = source["error"];
	        this.inputs = this.convertValues(source["inputs"], Input);
	        this.extra = this.convertValues(source["extra"], Input);
	        this.outputs = this.convertValues(source["outputs"], Output);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Activity {
	    key: string;
	    label: string;
	    icon_path: string;
	    description: string;
	    namespace?: string;
	    method?: string;
	    isPseudo?: boolean;
	    parameter_define?: ParameterDefine;
	    component?: string;
	    children?: Activity[];
	
	    static createFrom(source: any = {}) {
	        return new Activity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.icon_path = source["icon_path"];
	        this.description = source["description"];
	        this.namespace = source["namespace"];
	        this.method = source["method"];
	        this.isPseudo = source["isPseudo"];
	        this.parameter_define = this.convertValues(source["parameter_define"], ParameterDefine);
	        this.component = source["component"];
	        this.children = this.convertValues(source["children"], Activity);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

