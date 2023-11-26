export namespace models {
	
	export class Project {
	    id: number[];
	    name: string;
	    path: string;
	    icon: string;
	    description: string;
	    // Go type: time
	    create_ts: any;
	    // Go type: time
	    update_ts: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.icon = source["icon"];
	        this.description = source["description"];
	        this.create_ts = this.convertValues(source["create_ts"], null);
	        this.update_ts = this.convertValues(source["update_ts"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	
	export class Input {
	    key: string;
	    label: string;
	    category?: string;
	    type: string;
	    default_value: string;
	    required: boolean;
	    options: any;
	    editor_type: string;
	    show_if?: string;
	
	    static createFrom(source: any = {}) {
	        return new Input(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
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
	    error_deal?: boolean;
	    inputs?: Input[];
	    outputs?: Output[];
	
	    static createFrom(source: any = {}) {
	        return new ParameterDefine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.error_deal = source["error_deal"];
	        this.inputs = this.convertValues(source["inputs"], Input);
	        this.outputs = this.convertValues(source["outputs"], Output);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class Activitiy {
	    key: string;
	    label: string;
	    icon_path: string;
	    description: string;
	    namespace?: string;
	    method?: string;
	    isPseudo?: boolean;
	    parameter_define?: ParameterDefine;
	    component?: string;
	    children?: Activitiy[];
	
	    static createFrom(source: any = {}) {
	        return new Activitiy(source);
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
	        this.children = this.convertValues(source["children"], Activitiy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

