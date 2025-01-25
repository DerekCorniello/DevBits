/* Created with https://stirlingmarketinggroup.github.io/go2ts/*/

export interface UserProps {
    username: string;
    bio: string;
    links: string[];
    created_on: string;
    picture: string;
}

export interface ProjectProps {
    id: number;
    owner: number;
    name: string;
    description: string;
    status: number;
    likes: number;
    tags: string[];
    links: string[];
    creation_date: string;
}

export interface PostProps {
    id: number;
    user: number;
    project: number;
    likes: number;
    content: string;
    comments: number[];
    created_on: string;
}

export interface CommentProps {
    id: number;
    user: number;
    post: number;
    likes: number;
    parent_comment: number;
    created_on: string;
    content: string;
}

export interface ErrorResponseProps {
    error: string;
    message: string;
}