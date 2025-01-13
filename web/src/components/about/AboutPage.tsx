import "./AboutPage.css"

function AboutPage() {
    return (
        <div>
            This is a passion project, meant to be a playground for learning Full-Stack development<br />
            and hopefully one day something useful on its own.<br />
            <br />
            It is using the following technologies:<br />
             - GO lang for backend API services<br />
             - ReactJS, webpack and materialUI for frontend<br />
             - Docker for containerisation and GCP for infrastructure (Cloud Run, Cloud SQL)<br />
             - Google auth for signup/login<br />
             <br />
            Achitecture view:
            <img src="/high-level-arch.png" className="high-level-arch" alt="high level arch" />
        </div>
    )
}

export default AboutPage