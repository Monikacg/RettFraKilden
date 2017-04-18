%% Plot of states matrix
%load states_measured.mat



hold off

time_plot_102 = states_measured(1,:);
travel_plot_102 = states_measured(2,:);
travel_rate_plot_102 = states_measured(3,:);
pitch_plot_102 = states_measured(4,:);
pitch_rate_plot_102 = states_measured(5,:);
elevation_plot_102 = states_measured(6,:);
elevation_rate_plot_102 = states_measured(7,:);

subplot(711)
stairs(u_ts_padding(:,1),u_ts_padding(:,2)),grid
ylabel('u [rad]')
subplot(712)
plot(time_plot_102,travel_plot_102,'m',time_plot_102,travel_plot_102,'m','LineWidth',2),grid
ylabel('Travel [rad]')
subplot(713)
plot(time_plot_102,travel_rate_plot_102,'m',time_plot_102,travel_rate_plot_102','m','LineWidth',2),grid
ylabel('Travel rate [rad/s]')
subplot(714)
plot(time_plot_102,pitch_plot_102,'m',time_plot_102,pitch_plot_102','m','LineWidth',2),grid
ylabel('Pitch [rad]')
subplot(715)
plot(time_plot_102,pitch_rate_plot_102,'m',time_plot_102,pitch_rate_plot_102','m','LineWidth',2),grid
ylabel('Pitch rate [rad/s]')
subplot(716)
plot(time_plot_102,elevation_plot_102,'m',time_plot_102,elevation_plot_102','m','LineWidth',2),grid
ylabel('Elevation [rad]')
subplot(717)
plot(time_plot_102,elevation_rate_plot_102,'m',time_plot_102,elevation_rate_plot_102','m','LineWidth',2),grid
xlabel('tid [s]'),ylabel('Elevation rate [rad/s]')