export namespace DAO {
	
	export class Settings {
	    prefix: string;
	    proxy: string;
	    cookie: string;
	    r_18: boolean;
	    downloadposition: string;
	    likelimit: number;
	    retry429: number;
	    downloadinterval: number;
	    retryinterval: number;
	    differauthor: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prefix = source["prefix"];
	        this.proxy = source["proxy"];
	        this.cookie = source["cookie"];
	        this.r_18 = source["r_18"];
	        this.downloadposition = source["downloadposition"];
	        this.likelimit = source["likelimit"];
	        this.retry429 = source["retry429"];
	        this.downloadinterval = source["downloadinterval"];
	        this.retryinterval = source["retryinterval"];
	        this.differauthor = source["differauthor"];
	    }
	}

}

