export namespace models {
	
	export class Project {
	    id: number[];
	    name: string;
	    path: string;
	    icon: string;
	    description: string;
	    // Go type: JsonTime
	    create_ts: any;
	    // Go type: JsonTime
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

