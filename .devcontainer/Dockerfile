FROM library/ros:jazzy

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
	--mount=type=cache,target=/var/lib/apt,sharing=locked \
	rm -rf /etc/apt/apt.conf.d/docker-clean \
	&& apt update \
	&& apt-get install --yes --no-install-recommends \
		ros-jazzy-rmw-cyclonedds-cpp \
		zsh

COPY --from=library/golang:1.26 /usr/local/go /usr/local/go

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

ARG USERNAME=hypnos
RUN groupmod -n ${USERNAME} ubuntu \
	&& usermod -l ${USERNAME} ubuntu \
	&& usermod -d /home/${USERNAME} -m ${USERNAME} \
	&& echo "${USERNAME}:${USERNAME}" | chpasswd \
	&& echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > "/etc/sudoers.d/${USERNAME}"
	


ENV ROS_ROOT=/opt/ros/${ROS_DISTRO} \
	PYTHON_VERSION=3.12

ENV PATH=${ROS_ROOT}/bin:${PATH} \
	ROS_VERSION=2 \
	ROS_PYTHON_VERSION=3 \
	ROS_AUTOMATIC_DISCOVERY_RANGE=SUBNET \
	AMENT_PREFIX_PATH=${ROS_ROOT} \
	PYTHONPATH=${ROS_ROOT}/lib/python${PYTHON_VERSION}/site-packages:/usr/lib/python3/dist-packages \
	RMW_IMPLEMENTATION=rmw_cyclonedds_cpp



USER ${USERNAME}

RUN mkdir -p \
		~/go/pkg/mod \
		~/.cache/go-build

# Comment out the following lines if you don't want to use Zsh.
RUN curl -fsSL https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh | zsh || true \
	&& git clone https://github.com/zsh-users/zsh-autosuggestions         ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions     \
	&& git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting \
	&& git clone https://github.com/zsh-users/zsh-completions             ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-completions
COPY .zshrc /home/${USERNAME}
ENV SHELL=/bin/zsh
