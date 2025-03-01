function Navbar() {
  return (
    <div className="sticky top-0 w-full bg-neutral-200">
      <div className="py-3 max-w-7xl mx-auto border-b border-neutral-300 flex items-center gap-10">
        <img
          className="h-8"
          src="https://probo.in/_next/image?url=https%3A%2F%2Fd39axbyagw7ipf.cloudfront.net%2Fimages%2Flogo%2Flogo.webp&w=128&q=75"
          alt="logo"
        />
        <div className="flex items-center w-full justify-between">
          <div className="flex gap-10 text-xs tracking-widest">
            <a>Trading</a>
            <a>Team 11</a>
            <a>Read</a>
            <a>Cares</a>
            <a>Careers</a>
          </div>
          <div className="flex gap-5">
            <button className="bg-neutral-50 font-bold px-8 py-2 rounded text-sm border border-neutral-300">
              Download App
            </button>
            <button className="bg-neutral-950 text-neutral-50 font-medium px-8 py-2 rounded text-sm">
              Trade Online
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Navbar;
