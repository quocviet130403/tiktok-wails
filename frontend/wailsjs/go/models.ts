export namespace service {
	
	export class ProfileDouyin {
	    id: number;
	    nickname: string;
	    url: string;
	    last_video_reup: any;
	    retry_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ProfileDouyin(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.nickname = source["nickname"];
	        this.url = source["url"];
	        this.last_video_reup = source["last_video_reup"];
	        this.retry_count = source["retry_count"];
	    }
	}
	export class Profiles {
	    id: number;
	    name: string;
	    hashtag: string;
	    first_comment: string;
	    is_authenticated: boolean;
	    proxy_ip: string;
	    proxy_port: string;
	
	    static createFrom(source: any = {}) {
	        return new Profiles(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.hashtag = source["hashtag"];
	        this.first_comment = source["first_comment"];
	        this.is_authenticated = source["is_authenticated"];
	        this.proxy_ip = source["proxy_ip"];
	        this.proxy_port = source["proxy_port"];
	    }
	}
	export class Video {
	    ID: number;
	    Title: string;
	    VideoURL: string;
	    ThumbnailURL: string;
	    Duration: number;
	    LikeCount: number;
	    ProfileDouyinID: number;
	    Status: string;
	    IsDeleted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Title = source["Title"];
	        this.VideoURL = source["VideoURL"];
	        this.ThumbnailURL = source["ThumbnailURL"];
	        this.Duration = source["Duration"];
	        this.LikeCount = source["LikeCount"];
	        this.ProfileDouyinID = source["ProfileDouyinID"];
	        this.Status = source["Status"];
	        this.IsDeleted = source["IsDeleted"];
	    }
	}

}

