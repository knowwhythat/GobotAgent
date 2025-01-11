export namespace frontend {
	
	export class FileFilter {
	    DisplayName: string;
	    Pattern: string;
	
	    static createFrom(source: any = {}) {
	        return new FileFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DisplayName = source["DisplayName"];
	        this.Pattern = source["Pattern"];
	    }
	}

}

export namespace models {
	
	export class SettingData {
	    browserPath: string;
	    apiBaseUrl: string;
	    apiKey: string;
	    modelName: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browserPath = source["browserPath"];
	        this.apiBaseUrl = source["apiBaseUrl"];
	        this.apiKey = source["apiKey"];
	        this.modelName = source["modelName"];
	    }
	}

}

