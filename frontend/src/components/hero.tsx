function Hero() {
  return (
    <div className="h-screen">
      <div className="bg-neutral-200 h-1/2">
        <div className="bg-neutral-200 flex justify-between items-end max-w-7xl mx-auto h-full">
          <div className="bg-neutral-200 h-full w-full p-10">
            <h1 className="text-7xl mt-28 font-light tracking-wider font-sans">
              India's Leading
            </h1>
            <h2 className="text-6xl font-extralight tracking-normal mt-3 font-sans">
              Online Skill Gaming
            </h2>
            <h2 className="text-6xl font-extralight tracking-normal mt-3 font-sans">
              Platform
            </h2>

            <p className="text-xl tracking-wider font-medium text-neutral-500 mt-10">
              Sports, Entertainment, Economy or Finance.
            </p>

            <div className="mt-20 flex gap-5">
              <button className="bg-neutral-50 font-bold px-8 py-2 rounded text-sm border border-neutral-300">
                Download App
              </button>
              <button className="bg-neutral-950 text-neutral-50 font-medium px-8 py-2 rounded text-sm">
                Trade Online
              </button>
            </div>
          </div>
          <img
            src="https://probo.in/_next/image?url=https%3A%2F%2Fd39axbyagw7ipf.cloudfront.net%2Fimages%2Fhome%2Fheader%2Fheader-23012025.webp&w=640&q=75"
            className="max-w-[35rem] min-h-[37rem]"
          />
        </div>
      </div>
      <div className="bg-neutral-800 h-1/2">
        <div className="bg-neutral-800 grid grid-cols-2 max-w-7xl mx-auto">
          <div>
            <div>
              <span>Samachar</span>
              <span>Vichar</span>
              <span>Vyapaar</span>
            </div>
          </div>
          <div className="flex justify-center items-center p-4">
            <div className="border-2 border-neutral-700 p-2 rounded-[2.5rem]">
              <video
                autoPlay
                loop
                muted
                className="max-h-[37rem] bg-amber-50 rounded-4xl"
                src="https://d39axbyagw7ipf.cloudfront.net/videos/info-video.mp4"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Hero;
