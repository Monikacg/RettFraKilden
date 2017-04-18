%% Plot of states matrix
%load states.mat

hold off

time_plot_104 = states_measured(1,:);
travel_plot_104 = states_measured(2,:);
travel_rate_plot_104 = states_measured(3,:);
pitch_plot_104 = states_measured(4,:);
pitch_rate_plot_104 = states_measured(5,:);
elevation_plot_104 = states_measured(6,:);
elevation_rate_plot_104 = states_measured(7,:);

subplot(811)
stairs(u1_andtime(:,1),u1_andtime(:,2)),grid
ylabel('p_r_e_f')
subplot(812)
stairs(u2_andtime(:,1),u2_andtime(:,2)),grid
ylabel('e_r_e_f')
subplot(813)
plot(time_plot_104,travel_plot_104,'m',time_plot_104,travel_plot_104,'m','LineWidth',2),grid
ylabel('Travel')
subplot(814)
plot(time_plot_104,travel_rate_plot_104,'m',time_plot_104,travel_rate_plot_104','m','LineWidth',2),grid
ylabel('Travel rate')
subplot(815)
plot(time_plot_104,pitch_plot_104,'m',time_plot_104,pitch_plot_104','m','LineWidth',2),grid
ylabel('Pitch')
subplot(816)
plot(time_plot_104,pitch_rate_plot_104,'m',time_plot_104,pitch_rate_plot_104','m','LineWidth',2),grid
ylabel('Pitch rate')
subplot(817)
plot(time_plot_104,elevation_plot_104,'m',time_plot_104,elevation_plot_104','m','LineWidth',2),grid
ylabel('Elevation')
subplot(818)
plot(time_plot_104,elevation_rate_plot_104,'m',time_plot_104,elevation_rate_plot_104','m','LineWidth',2),grid
xlabel('tid (s)'),ylabel('Elevation rate')