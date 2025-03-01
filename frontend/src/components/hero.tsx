function Hero() {
  return (
    <div className="">
      <div className="bg-neutral-200">
        <div className="bg-neutral-200 flex justify-between items-center max-w-7xl mx-auto h-full">
          <div className="bg-neutral-200 h-full w-full p-10">
            <h1 className="text-7xl font-light tracking-wider font-sans">
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
          <div className="w-[80rem] h-[40rem] flex items-end justify-center">
            <img
              src="https://probo.in/_next/image?url=https%3A%2F%2Fd39axbyagw7ipf.cloudfront.net%2Fimages%2Fhome%2Fheader%2Fheader-23012025.webp&w=640&q=75"
              className="w-auto h-auto max-w-full max-h-full object-contain"
            />
          </div>
        </div>
      </div>
      <div className="bg-neutral-800 h-1/2">
        <div className="bg-neutral-800 grid grid-cols-2 max-w-7xl mx-auto">
          <div>
            <div className="bg-neutral-800 h-full flex flex-col justify-center px-4">
              <div className="flex gap-10 text-4xl tracking-wider font-medium mb-8">
                <span className="text-neutral-100">Samachar</span>
                <span className="text-neutral-500">Vichaar</span>
                <span className="text-neutral-500">Vyapaar</span>
              </div>
              <p className="text-2xl tracking-widest text-neutral-200">
                Lorem ipsum dolor sit amet, consectetur adipisicing elit. Quas
                esse, repudiandae alias laborum at impedit quo eveniet ea
              </p>
            </div>
          </div>
          <div className="flex justify-center items-center p-10">
            <div className="border-2 border-neutral-700 p-2 rounded-[2.5rem]">
              <video
                autoPlay
                loop
                muted
                className="max-h-[35rem] bg-amber-50 rounded-4xl"
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
